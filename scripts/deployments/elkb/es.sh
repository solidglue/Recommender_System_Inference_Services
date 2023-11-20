#!/bin/bash

$ docker run -d --name elasticsearch  docker.elastic.co/elasticsearch/elasticsearch:5.4.0


# 
$ docker exec -it elasticsearch bash

# edit config file
$ vim config/elasticsearch.yml
cluster.name: "docker-cluster"
node.name: node-101
path.data: /data/esdata	#data dir
path.logs: /data/eslogs	#logs dir
network.host: 0.0.0.0
http.cors.enabled: true
http.cors.allow-origin: "*"
xpack.security.enabled: false
# minimum_master_nodes need to be explicitly set when bound on a public IP
# set to 1 to allow single node clusters
# Details: https://github.com/elastic/elasticsearch/pull/17288
discovery.zen.ping.unicast.hosts: ["192.168.9.101", "192.168.9.102","192.168.9.103"]
discovery.zen.minimum_master_nodes: 1
http.port: 9200 	#es port


docker commit -a "add config" -m "dev" a404c6c174a2  es:latest
docker run -d --name es -p 9200:9200 -p 9300:9300   -e "discovery.type=single-node" es:latest
$ curl 'http://localhost:9200/_nodes/http?pretty'


docker run -p 9100:9100 mobz/elasticsearch-head:5

