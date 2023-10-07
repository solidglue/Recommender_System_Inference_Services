package cores

import (
	"encoding/json"
	"infer-microservices/common"
	faiss_index "infer-microservices/common/faiss_gogofaster"
	"infer-microservices/cores/service_config"
	"infer-microservices/utils"
	"infer-microservices/utils/logs"
	"net/http"
	"sync"
	"time"

	"github.com/allegro/bigcache"
)

var bigCacheConfRankResult bigcache.Config

type DeepFM struct {
	userId        string
	itemList      []string
	serviceConfig *service_config.ServiceConfig
}

func init() {
	bigCacheConfRankResult = bigcache.Config{
		Shards:             shards,
		LifeWindow:         lifeWindowS * time.Minute,
		CleanWindow:        cleanWindowS * time.Minute,
		MaxEntriesInWindow: maxEntriesInWindow,
		MaxEntrySize:       maxEntrySize,
		Verbose:            verbose,
		HardMaxCacheSize:   hardMaxCacheSize,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}
}

// userid
func (d *DeepFM) SetUserId(userId string) {
	d.userId = userId
}

func (d *DeepFM) getUserId() string {
	return d.userId
}

// itemList
func (d *DeepFM) SetItemList(itemList []string) {
	d.itemList = itemList
}

func (d *DeepFM) getItemList() []string {
	return d.itemList
}

// serviceConfig *service_config.ServiceConfig
func (d *DeepFM) SetServiceConfig(serviceConfig *service_config.ServiceConfig) {
	d.serviceConfig = serviceConfig
}

func (d *DeepFM) getServiceConfig() *service_config.ServiceConfig {
	return d.serviceConfig
}

func (d *DeepFM) RankInferSkywalking(r *http.Request) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)

	tensorName := "scores"
	cacheKeyPrefix := d.userId + d.serviceConfig.GetServiceId() + "_rankResult"

	//set cache
	bigCache, err := bigcache.NewBigCache(bigCacheConfRankResult)
	if err != nil {
		logs.Error(err)
	}

	// get features from cache.
	if lifeWindowS > 0 {
		exampleDataBytes, _ := bigCache.Get(cacheKeyPrefix)
		err = json.Unmarshal(exampleDataBytes, &response)
		if err != nil {
			return nil, err
		}
		return response, nil
	}

	//get infer samples.
	spanUnionEmFv, _, err := common.Tracer.CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}
	spanUnionEmFv.SetOperationName("get rank infer examples func")
	spanUnionEmFv.Log(time.Now())
	examples, err := d.getInferExampleFeatures()
	if err != nil {
		return nil, err
	}
	spanUnionEmFv.Log(time.Now())
	spanUnionEmFv.End()

	// get rank scores from tfserving model.
	items := make([]string, 0)
	scores := make([]float32, 0)
	rankResult := make([]*faiss_index.ItemInfo, 0)
	spanUnionEmFv, _, err = common.Tracer.CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}
	spanUnionEmFv.SetOperationName("get rank scores func")
	spanUnionEmFv.Log(time.Now())
	items_, scores_, err := d.getServiceConfig().GetModelClient().RankPredict(examples, tensorName)
	if err != nil {
		return nil, err
	} else {
		items = *items_
		scores = *scores_
	}
	spanUnionEmFv.Log(time.Now())
	spanUnionEmFv.End()

	//build rank result whith tfserving.ItemInfo
	for idx := 0; idx < len(items); idx++ {
		itemInfo := &faiss_index.ItemInfo{
			ItemId: items[idx],
			Score:  scores[idx],
		}
		rankResult = append(rankResult, itemInfo)
	}

	//format result.
	spanUnionEmOut, _, err := common.Tracer.CreateLocalSpan(r.Context())
	if err != nil {
		return nil, err
	}
	spanUnionEmOut.SetOperationName("get rank result func")
	spanUnionEmOut.Log(time.Now())
	rankRst, err := d.recallResultFmt(&rankResult)
	if err != nil {
		return nil, err
	}
	spanUnionEmOut.Log(time.Now())
	spanUnionEmOut.End()

	response["data"] = *rankRst
	if lifeWindowS > 0 {
		bigCache.Set(cacheKeyPrefix, []byte(utils.Struct2Json(response)))
	}

	return response, nil
}

func (d *DeepFM) RankInferNoSkywalking(r *http.Request) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
	tensorName := "scores"
	cacheKeyPrefix := d.getUserId() + d.serviceConfig.GetServiceId() + "_rankResult"

	//set cache
	bigCache, err := bigcache.NewBigCache(bigCacheConfRankResult)
	if err != nil {
		logs.Error(err)
	}

	// get features from cache.
	if lifeWindowS > 0 {
		exampleDataBytes, _ := bigCache.Get(cacheKeyPrefix)
		err = json.Unmarshal(exampleDataBytes, &response)
		if err != nil {
			return nil, err
		}
		return response, nil

	}

	//get infer samples.
	examples, err := d.getInferExampleFeatures()
	if err != nil {
		return nil, err
	}

	// get rank scores from tfserving model.
	items := make([]string, 0)
	scores := make([]float32, 0)
	rankResult := make([]*faiss_index.ItemInfo, 0)
	items_, scores_, err := d.getServiceConfig().GetModelClient().RankPredict(examples, tensorName)
	if err != nil {
		logs.Error(err)
		return nil, err
	} else {
		items = *items_
		scores = *scores_
	}

	//build rank result whith tfserving.ItemInfo
	for idx := 0; idx < len(items); idx++ {
		itemInfo := &faiss_index.ItemInfo{
			ItemId: items[idx],
			Score:  scores[idx],
		}
		rankResult = append(rankResult, itemInfo)
	}

	//format result.
	rankRst, err := d.recallResultFmt(&rankResult)
	if err != nil {
		return nil, err
	}

	response["data"] = *rankRst

	if lifeWindowS > 0 {
		bigCache.Set(cacheKeyPrefix, []byte(utils.Struct2Json(response)))
	}

	return response, nil
}

func (d *DeepFM) recallResultFmt(rankResult *[]*faiss_index.ItemInfo) (*[]map[string]interface{}, error) {

	recall := make([]map[string]interface{}, 0)
	recallTmp := make(chan map[string]interface{}, len(*rankResult)) // 20221011
	var wg sync.WaitGroup

	for idx := 0; idx < len(*rankResult); idx++ {
		rawCell := (*rankResult)[idx]
		wg.Add(1)
		go func(raw_cell_ *faiss_index.ItemInfo) {
			defer wg.Done()
			returnCell := make(map[string]interface{})
			returnCell["itemid"] = raw_cell_.ItemId
			returnCell["score"] = utils.FloatRound(raw_cell_.Score, 4)
			recallTmp <- returnCell
		}(rawCell)

	}
	wg.Wait()

	for idx := 0; idx < len(*rankResult); idx++ {
		returnCellTmp := <-recallTmp
		recall = append(recall, returnCellTmp)
	}
	close(recallTmp)

	return &recall, nil
}
