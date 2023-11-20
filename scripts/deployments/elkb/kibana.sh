#/bin/bash


vi config/kibana.yml
server.port: 5601		#kibana port
server.host: "127.0.0.1"		#kibana ip
server.name: "infer-kibana"		#kibana name
elasticsearch.url: http://127.0.0.1:9200 		#es addr

docker run --name kibana -e ELASTICSEARCH_URL=http://127.0.0.1:9200 -p 5601:5601 -d kibana:5.6.9
