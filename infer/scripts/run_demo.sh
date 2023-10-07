#!/bin/bash

today=`date -d "0 days ago" "+%Y%m%d"`
go_project=infer
WORK_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
LOGPATH=${WORK_PATH}/logs

echo $WORK_PATH
echo $LOGPATH

#conf file
service_start_file="./conf/server_start_file.json"

#dubbo
dubbo_server_conf1="./conf/dubbogo_server_conf1.yml"
dubbo_server_conf2="./conf/dubbogo_server_conf2.yml"

#nacos
nacos_ip="10.10.10.11"
nacos_port=1022
nacos_username="nacos"
nacos_password="nacos"
nacos_loglevel="error"

#hystrix
hystrix_timeoutMS=100
hystrix_lowerRecallNum=100
hystrix_lowerRankNum=20

#skywalking
skywalking_ip="10.10.10.10"
skywalking_port=1020
skywalking_whetheropen=true
skywalking_servername="infer"

#redis
redis_password="#67fag!@#"

#port
rest_server_port=8020
grpc_server_port=8021
max_cpu_num=20

#bigcache
bigcahe_lifeWindowS=300
bigcache_cleanWindowS=120
bigcache_hardMaxCacheSize=409600

#tensorflow
tfserving_timeoutms=100

#logs
log_level="error"

function status(){
   ps -ef | grep ${go_project}
}

function start(){

    echo "starting servers ..."
    cd ${WORK_PATH}/
    dubbo_conf=$1
    
    nohup ${WORK_PATH}/${go_project} \
    
        #start file
        --service_start_file
        
        #dubbo
        --dubbo_server_conf=${dubbo_conf} \
        
        #nacos
        --nacos_ip="10.10.10.11" \
        --nacos_port=1022
        --nacos_username="nacos" \
        --nacos_password="nacos" \
        --nacos_loglevel="error" \
        
        #hystrix
        --hystrix_timeoutMS=100  \
        --hystrix_lowerRecallNum=100  \
        --hystrix_lowerRankNum=20 \
        
        #skywalking
        --skywalking_ip="10.10.10.10" \
        --skywalking_port=1020 \
        --skywalking_whetheropen=true \
        --skywalking_servername="infer" \
        
        #redis
        --redis_password="#67fag!@#" \
        
        #server port
        --rest_server_port=8020 \
        --grpc_server_port=8021 \
        --max_cpu_num
        
        #bigcache
        --bigcahe_lifeWindowS=300 \
        --bigcache_cleanWindowS=120 \
        --bigcache_hardMaxCacheSize=409600 \
        
        #tensorflow
        --tfserving_timeoutms=100 \
        
        #logs
        --log_level="error"   > ${LOGPATH}/${go_project}_${today}.log 2>&1 &

}

function stop(){
    ps -ef | grep ${go_project} | grep -v grep | awk '{print $2}' | xargs kill -9
}

function main(){
    case $1 in
        "status")
            status
        ;;

        "start")  
            #start 2 services.
            start ${dubbo_server_conf1}
            start ${dubbo_server_conf2}
            status
        ;;

        "stop")
            status
            stop
        ;;
		
		*)
		    echo "run args err, please input (status | start | stop) "	
	esac		
}

main

