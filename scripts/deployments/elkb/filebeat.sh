#!/bin/bash


wget https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-5.4.3-linux-x86_64.tar.gz
tar -zxvf filebeat-5.4.3-linux-x86_64.tar.gz
mv filebeat-5.4.3-linux-x86_64 filebeat


#config multi log files
filebeat.inputs:

- type: log
  enabled: true
  paths:
    - /var/logs/info/*
  tags: ["infer-info-1"]

- type: log
  enabled: true
  paths:
    - /var/logs/error/*
  tags: ["infer-info-2"]

output.logstash:
  hosts: ["localhost:5044"]
#output.elasticsearch:
#  # Array of hosts to connect to.
#  hosts: ["192.168.81.129:9200"]
#  indices:
#    - index: "infer-info-1-%{+yyyy.MM.dd}"
#      when.contains:
#        tags: "infer-info-1"
#    - index: "infer-info-2-%{+yyyy.MM.dd}"
#      when.contains:
#        tags: "infer-info-2"

./filebeat  -e  -c client.yml
