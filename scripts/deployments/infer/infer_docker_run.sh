#/bin/bash

#6020 rest port ,6021 grpc port
docker run -d --name infer -p 22000:6020 -p 22001:6021 --restart=always  recsys-go-infer:v1.0.0
curl localhost:22000