#!/bin/bash

curl http://10.124.20.5:8650/infer -X POST -d 'data={"modelId":"deepfm-1|infer","userId":"112223333444555","itemIdList":["111","222","333","444"]}' 

