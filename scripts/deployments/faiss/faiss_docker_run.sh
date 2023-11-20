#/bin/bash

docker run -d --name infer -p 21000:5000 --restart=always  faiss:v1.0.0
curl localhost:21000
