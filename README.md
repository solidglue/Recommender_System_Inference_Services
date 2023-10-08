# RecommenderSystems Inference Microservices
    DeepLearning Based Recommender Systems Infer Microsevices for Golang.

## Dependent Components    
The model inference microservices based on deep learning mainly uses the following components:
| Type | Component | Description |
| --- | --- | --- |
| Data | Hive / Spark | ETL millions users's behavior data and then build the feature data warehouse. |
|| Redis |  Save the training samples in TFRecord format and store them in Redis Cluster. |
| Models | TensorFlow | Training deep learning recall / rank model , alse you can use other deep learning framework ,but need save models as *.pb format. |
|| TensorFlow Serving | Deploy models and provide a grpc service. |
||FAISS | Quick search thousands items from millions items. |
| Microservices | Nacos | Manage config files and services. |
|| Dubbogo | Build dubbo protocol RPC services and register them to Nacos. |
|| Hystrix | How to distribute traffic during peak traffic (Latency and Fault Tolerance). |
|| Skywalking | Record the time spent on each request. |
|| Kubernetes  | Manage dockers and monitor the resource consumption of each service, such as memory and CPUs. |
||  Nginx、Apisix | API gateway. |



## Architecture

| Type | Component | Description |
| --- | --- | --- |
| Samples | [Recall Samples](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/blob/master/cores/dssm_samples.go) | Search recall TFRcords format samples from redis cluster, such as dssm model. |
|| [Rank Samples](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/blob/master/cores/deepfm_samples.go) |  Search recall TFRcords format samples from redis cluster, such as deepfm model. |
| Recall | [Get user/item vector](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/tree/master/cores/model) | Search user's embedding vector from recall model which deployed by tfservig(grpc sevice) , input data is recall samples. |
|  | [Search index](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/tree/master/cores/faiss) | Quick search thousands items from faiss index (millions items) service(grpc sevice) . |
| Rank | [Rank](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/blob/master/cores/rank_infer_deepfm.go)  | Rank input items by rank model  which deployed by tfservig(grpc sevice) . |
| Services | [Config Loader](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/tree/master/cores/service_config) | Sparse service's start config from Naocs, such as grpc info 、 redis info and index info. |
|  | [Register Services](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/blob/master/apis/dubbo/server/dubbo_server.go) | Register services to Nacos. |
|  | [Update Services](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/tree/master/cores/nacos_config) | Update services when nacos config files have changed, such as grpc info 、 redis info or index info. |
|APIS| [Dubbo](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/tree/master/apis/dubbo) |Provide Dubbo protocol APIs. |
|| [gRPC](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/tree/master/apis/grpc) |Provide gRPC protocol APIs. |
|| [REST](https://github.com/beachdogs/RecommenderSystems-Inference-Microservices/tree/master/apis/rest) |Provide Http protocol APIs. |



## Service Deploy
    Docker


    
