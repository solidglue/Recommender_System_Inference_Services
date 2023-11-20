# start script in docker.

#!/bin/bash

today=`date -d "0 days ago" "+%Y%m%d"`
go_project=recsys-go-infer
WORK_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
LOGPATH=${WORK_PATH}/logs

echo $WORK_PATH
echo $LOGPATH

#conf file
service_start_file="./configs/server_start_config.json"

#dubbo
dubbo_server_conf1="./configs/dubbogo_server-1_config.yml"
dubbo_server_conf2="./configs/dubbogo_server-2_config.yml"

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
rest_server_port=6020
grpc_server_port=6021
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
        --nacos_ip=${nacos_ip} \
        --nacos_port=${nacos_port}
        --nacos_username=${nacos_username} \
        --nacos_password=${nacos_password} \
        --nacos_loglevel=${nacos_loglevel} \
        
        #hystrix
        --hystrix_timeoutMS=${hystrix_timeoutMS}  \
        --hystrix_lowerRecallNum=${hystrix_lowerRecallNum}  \
        --hystrix_lowerRankNum=${hystrix_lowerRankNum} \
        
        #skywalking
        --skywalking_ip=${skywalking_ip} \
        --skywalking_port=${skywalking_port} \
        --skywalking_whetheropen=${skywalking_whetheropen} \
        --skywalking_servername=${skywalking_servername} \
        
        #redis
        --redis_password=${redis_password} \
        
        #server port
        --rest_server_port=${rest_server_port} \
        --grpc_server_port=${grpc_server_port} \
        --max_cpu_num=${max_cpu_num}   \
        
        #bigcache
        --bigcahe_lifeWindowS=${bigcahe_lifeWindowS} \
        --bigcache_cleanWindowS=${bigcache_cleanWindowS} \
        --bigcache_hardMaxCacheSize=${bigcache_hardMaxCacheSize} \
        
        #tensorflow
        --tfserving_timeoutms=${tfserving_timeoutms} \
        
        #logs
        --log_level=${log_level}   > ${LOGPATH}/${go_project}_${today}.log 2>&1 &

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

