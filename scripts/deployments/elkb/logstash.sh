#/bin/bash

wget https://artifacts.elastic.co/downloads/logstash/logstash-5.4.3.tar.gz
tar -zxvf logstash-5.4.3.tar.gz


vi logstash-5.4.3/client.conf
input {  # input config, use filebeat
    beats {
        port => 5044                 #filebeat and logstash deploy at the same servers
        codec => "json"
    }
}

#multi log files and multi es index
filter{
    if "infer-info-1"  in [tags] {
        grok {
            match => {
                "message"=>"" #filter depends on bussiness
            }
        }
    }

    if "infer-info-2"  in [tags] {
        grok {
            match=>{
                "message"=>"" #filter depends on bussiness
            }
        }
    }

    mutate {  #remove
        remove_field =>["message"]
        remove_field =>["beat"]
        remove_field =>["host"]
        #remove_field =>["fields"]
        remove_field =>["input"]
        remove_field =>["prospector"]
    }
}

output{
    if "infer-info-1"  in [tags] {
        elasticsearch {
            hosts => ["http://x.x.x.130:9200"]
            index => "infer-info-1-%{+yyyy.MM.dd}"
            #user => "elastic"
            #password => "changeme"
        }    
         stdout {codec => rubydebug}  
    }
    if  "infer-info-1"  in [tags] {
        elasticsearch {
            hosts => ["http://x.x.x.130:9200"]
            index => "infer-info-2-%{+yyyy.MM.dd}"
            #user => "elastic"
            #password => "changeme"
        }
        stdout {codec => rubydebug}
    }
}



bin/logstash  -f client.conf