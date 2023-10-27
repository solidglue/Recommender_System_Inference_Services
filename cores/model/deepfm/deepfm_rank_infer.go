package deepfm

import (
	"context"
	"encoding/json"
	"infer-microservices/common"
	faiss_index "infer-microservices/common/faiss_gogofaster"
	"infer-microservices/common/flags"
	framework "infer-microservices/common/tensorflow_gogofaster/core/framework"
	tfserving "infer-microservices/common/tfserving_gogofaster"
	"infer-microservices/cores/service_config_loader"
	"infer-microservices/utils"
	"infer-microservices/utils/logs"
	"net/http"
	"sync"
	"time"

	"github.com/allegro/bigcache"
	"github.com/gogo/protobuf/types"
)

var bigCacheConfRankResult bigcache.Config
var tfservingModelVersion int64
var tfservingTimeout int64

type DeepFM struct {
	userId        string
	retNum        int
	itemList      []string
	serviceConfig *service_config_loader.ServiceConfig
}

func init() {
	bigCacheConfRankResult = bigcache.Config{
		Shards:             shards1,
		LifeWindow:         lifeWindowS1 * time.Minute,
		CleanWindow:        cleanWindowS1 * time.Minute,
		MaxEntriesInWindow: maxEntriesInWindow1,
		MaxEntrySize:       maxEntrySize1,
		Verbose:            verbose1,
		HardMaxCacheSize:   hardMaxCacheSize1,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}

	flagFactory := flags.FlagFactory{}
	flagTensorflow := flagFactory.FlagTensorflowFactory()

	tfservingModelVersion = *flagTensorflow.GetTfservingModelVersion()
	tfservingTimeout = *flagTensorflow.GetTfservingTimeoutMs()
}

// userid
func (d *DeepFM) SetUserId(userId string) {
	d.userId = userId
}

func (d *DeepFM) getUserId() string {
	return d.userId
}

// retNum
func (d *DeepFM) SetRetNum(retNum int) {
	d.retNum = retNum
}

func (d *DeepFM) getRetNum() int {
	return d.retNum
}

// itemList
func (d *DeepFM) SetItemList(itemList []string) {
	d.itemList = itemList
}

func (d *DeepFM) getItemList() []string {
	return d.itemList
}

// serviceConfig *service_config.ServiceConfig
func (d *DeepFM) SetServiceConfig(serviceConfig *service_config_loader.ServiceConfig) {
	d.serviceConfig = serviceConfig
}

func (d *DeepFM) getServiceConfig() *service_config_loader.ServiceConfig {
	return d.serviceConfig
}

func (d *DeepFM) ModelInferSkywalking(r *http.Request) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)

	tensorName := "scores"
	cacheKeyPrefix := d.userId + d.serviceConfig.GetServiceId() + "_rankResult"

	//set cache
	bigCache, err := bigcache.NewBigCache(bigCacheConfRankResult)
	if err != nil {
		logs.Error(err)
	}

	// get features from cache.
	if lifeWindowS1 > 0 {
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
	examples, err := d.GetInferExampleFeatures()
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
	items_, scores_, err := d.rankPredict(examples, tensorName) // d.getServiceConfig().GetModelConfig().rankPredict(examples, tensorName)
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
	rankRst, err := d.rankResultFmt(&rankResult)
	if err != nil {
		return nil, err
	}
	spanUnionEmOut.Log(time.Now())
	spanUnionEmOut.End()

	response["data"] = *rankRst
	if lifeWindowS1 > 0 {
		bigCache.Set(cacheKeyPrefix, []byte(utils.Struct2Json(response)))
	}

	return response, nil
}

func (d *DeepFM) ModelInferNoSkywalking(r *http.Request) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
	tensorName := "scores"
	cacheKeyPrefix := d.getUserId() + d.serviceConfig.GetServiceId() + "_rankResult"

	//set cache
	bigCache, err := bigcache.NewBigCache(bigCacheConfRankResult)
	if err != nil {
		logs.Error(err)
	}

	// get features from cache.
	if lifeWindowS1 > 0 {
		exampleDataBytes, _ := bigCache.Get(cacheKeyPrefix)
		err = json.Unmarshal(exampleDataBytes, &response)
		if err != nil {
			return nil, err
		}
		return response, nil

	}

	//get infer samples.
	examples, err := d.GetInferExampleFeatures()
	if err != nil {
		return nil, err
	}

	// get rank scores from tfserving model.
	items := make([]string, 0)
	scores := make([]float32, 0)
	rankResult := make([]*faiss_index.ItemInfo, 0)
	items_, scores_, err := d.rankPredict(examples, tensorName) //d.getServiceConfig().GetModelConfig().rankPredict(examples, tensorName)
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
	rankRst, err := d.rankResultFmt(&rankResult)
	if err != nil {
		return nil, err
	}

	response["data"] = *rankRst

	if lifeWindowS1 > 0 {
		bigCache.Set(cacheKeyPrefix, []byte(utils.Struct2Json(response)))
	}

	return response, nil
}

func (d *DeepFM) rankResultFmt(rankResult *[]*faiss_index.ItemInfo) (*[]map[string]interface{}, error) {

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

// request rank scores from tfserving
func (d *DeepFM) rankPredict(examples common.ExampleFeatures, tensorName string) (*[]string, *[]float32, error) {

	userExamples := make([][]byte, 0)
	userContextExamples := make([][]byte, 0)
	itemExamples := make([][]byte, 0)
	items := make([]string, 0)

	userExamples = append(userExamples, *(examples.UserExampleFeatures.Buff))
	userContextExamples = append(userContextExamples, *(examples.UserContextExampleFeatures.Buff))

	for _, itemExample := range *examples.ItemSeqExampleFeatures {
		items = append(items, *(itemExample.Key))
		itemExamples = append(itemExamples, *(itemExample.Buff))
	}

	scores, err := d.requestTfservering(&userExamples, &userContextExamples, &itemExamples, tensorName)

	if err != nil {
		logs.Error(err)
		return nil, nil, err
	}

	return &items, scores, nil
}

func (d *DeepFM) requestTfservering(userExamples *[][]byte, userContextExamples *[][]byte, itemExamples *[][]byte, tensorName string) (*[]float32, error) {
	grpcConn, err := d.getServiceConfig().GetModelConfig().GetTfservingGrpcPool().Get()
	defer d.getServiceConfig().GetModelConfig().GetTfservingGrpcPool().Put(grpcConn)

	if err != nil {
		return nil, err
	}
	predictConfig := tfserving.NewPredictionServiceClient(grpcConn)

	version := &types.Int64Value{Value: tfservingModelVersion}
	predictRequest := &tfserving.PredictRequest{
		ModelSpec: &tfserving.ModelSpec{
			Name:    d.getServiceConfig().GetModelConfig().GetModelName(),
			Version: version,
		},
		Inputs: make(map[string]*framework.TensorProto),
	}

	//user examples
	tensorProtoUser := &framework.TensorProto{
		Dtype: framework.DataType_DT_STRING,
	}
	tensorProtoUser.TensorShape = &framework.TensorShapeProto{
		Dim: []*framework.TensorShapeProto_Dim{
			{
				Size_: int64(len(*userExamples)),
				Name:  "",
			},
		},
	}
	tensorProtoUser.StringVal = *userExamples
	predictRequest.Inputs["userExamples"] = tensorProtoUser

	//context examples, realtime
	tensorProtoUserContext := &framework.TensorProto{
		Dtype: framework.DataType_DT_STRING,
	}
	tensorProtoUserContext.TensorShape = &framework.TensorShapeProto{
		Dim: []*framework.TensorShapeProto_Dim{
			{
				Size_: int64(len(*userContextExamples)),
				Name:  "",
			},
		},
	}
	tensorProtoUserContext.StringVal = *userContextExamples
	predictRequest.Inputs["userContextExamples"] = tensorProtoUserContext

	//item examples
	tensorProtoItem := &framework.TensorProto{
		Dtype: framework.DataType_DT_STRING,
	}
	tensorProtoItem.TensorShape = &framework.TensorShapeProto{
		Dim: []*framework.TensorShapeProto_Dim{
			{
				Size_: int64(len(*itemExamples)),
				Name:  "",
			},
		},
	}
	tensorProtoItem.StringVal = *itemExamples
	predictRequest.Inputs["itemExamples"] = tensorProtoItem

	predictRequest.OutputFilter = []string{tensorName}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(tfservingTimeout)*time.Millisecond)
	defer cancel()

	predict, err := predictConfig.Predict(ctx, predictRequest)
	if err != nil {
		return nil, err
	}
	predictOut, _ := predict.Outputs[tensorName]

	return &predictOut.FloatVal, nil
}
