#!/bin/bash

curl http://10.124.20.5:8650/infer -X POST -d 'data={"modelId":"deepfm-1|modelInfer","userId":"13438935173","itemIdList":["111","222","333","444"]}' 

