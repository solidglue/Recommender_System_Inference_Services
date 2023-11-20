#/bin/bash

docker run -p 20888:8500 -p 20881:8501 \
--mount type=bind,source=/data/tensorflow/models/deepfm,target=/models/deepfm \
--mount type=bind,source=/data/tensorflow/config/deepfm/model.conf,target=/models/model.conf \
-t tfserving:v1.15.0  \
--model_config_file=/models/model.conf
--restart=always 
