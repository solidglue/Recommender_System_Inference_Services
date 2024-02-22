# RecommenderSystems Inference Microservices
    DeepLearning Based Recommender Systems Infer Microsevices for Golang.

## Dependent Components    
The model inference microservices based on deep learning mainly uses the following components:
| Type | Component | Description |
| --- | --- | --- |
| Data | Hive / Spark | ETL millions users's behavior data and then build the feature data warehouse. |
|| Redis |  Save the training samples in TFRecord format and store them in Redis Cluster. |
| Model | TensorFlow | Training deep learning recall / rank model , alse you can use other deep learning framework ,but need save models as *.pb format. |
|| TensorFlow Serving | Deploy models and provide a grpc service. |
||FAISS | Quick search thousands items from millions items. |
| Microservices | Nacos | Manage config files and services. |
|| Dubbo | Build dubbo protocol RPC services and register them to Nacos. |
|| Hystrix | How to distribute traffic during peak traffic (Latency and Fault Tolerance). |
|| Skywalking | Record the time spent on each request. |
|Deploy| Docker  | Docker containerization deployment services. |
|| Kubernetes  | Manage dockers and monitor the resource consumption of each service, such as memory and CPUs. |
||  Nginx、Apisix | API gateway. |



## Architecture
The core components of model inference microservices are as follows：
| Type | Component | Description |
| --- | --- | --- |
| Feature | [feature engineering](https://github.com/solidglue/RecommenderSystems-Inference-Microservices/tree/master/pkg/infer_samples/feature) | user offline、user realtime、user seq features and item features. |
| Sample | [recall/rank samples](https://github.com/solidglue/RecommenderSystems-Inference-Microservices/tree/master/pkg/infer_samples) | create TFRcords format samples. |
| Recall | [cf recall](https://github.com/solidglue/RecommenderSystems-Inference-Microservices/tree/master/pkg/infer_models/recall/cf) | user cf 、 item cf and swing.|
|  | [dssm recall](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/tree/master/pkg/infer_models/recall/dssm) | recall from dssm model and faiss index. |
|  | [rules recall](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/tree/master/pkg/infer_models/recall/rules_recall) | rules recall, such as hot items recall. |
|  | [cold start](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/tree/master/pkg/infer_models/recall/cold_start) | new users and new items cold start. |
| Rank | [pre_ranking](https://github.com/solidglue/RecommenderSystems-Inference-Microservices/tree/master/pkg/infer_model/pre_ranking)  |  thousands items pre_ranking after recall . |
|  | [ranking](https://github.com/solidglue/RecommenderSystems-Inference-Microservices/tree/master/pkg/infer_model/ranking)  | hundreds items ranking after pre_ranking. |
|  | [re_ranking](https://github.com/solidglue/RecommenderSystems-Inference-Microservices/tree/master/pkg/infer_model/re_ranking)  |  hundreds items re_ranking after ranking . |
| Services | [config loader](https://github.com/solidglue/RecommenderSystems-Inference-Microservices/tree/master/pkg/config_loader) | Sparse service's start config from Naocs, such as grpc info 、 redis info and index info. |
|  | [dubbo service](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/blob/master/pkg/infer_services/dubbo_service) | dubbo protocol service. |
|  | [gRPC service](https://github.com/solidglue/RecommenderSystems-Inference-Microservices/tree/master/pkg/infer_services/grpc_service) | grpc protocol service. |
|  | [rest service](https://github.com/solidglue/RecommenderSystems-Inference-Microservices/tree/master/pkg/infer_services/rest_service) | restful service. |
|APIs| [dubbo api](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/tree/master/api/dubbo_api) |Provide Dubbo protocol APIs. |
|| [gRPC api](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/tree/master/api/grpc_api) |provide gRPC protocol APIs. |
|| [rest api](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/tree/master/api/rest_api) |provide Http protocol APIs. |
|Web| [web](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/tree/master/pkg/web) |services manage and Service monitor page. |
|Deploy| [faiss](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/tree/master/scripts/deployments/faiss) |faiss index service deploy. |
|| [tfserving](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/tree/master/scripts/deployments/tfserving) |tensorflow model deploy. |
|| [infer](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/tree/master/scripts/deployments/infer) |recommend system infer deploy. |




## Services Deploy
    Docker
    Kubernetes 
    Nginx
    Apisix
