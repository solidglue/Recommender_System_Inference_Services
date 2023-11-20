#/bin/bash

docker build  -f /data/loki/the-infer/Dockerfile/tfserving/tfserving_docker_file  -t tfserving:v1.15.0 .