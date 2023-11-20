#/bin/bash

#build docker mir
docker build -f /data/loki/the-infer/Dockerfile/infer/infer_docker_file -t recsys-go-infer:v1.0.0 .