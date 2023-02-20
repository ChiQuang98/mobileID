#!/bin/bash
nohup go run main.go &

#command to get source
#tcpdump -i any -n -c 5 ip | awk '{ print gensub(/(.*)\..*/,"\\1","g",$3), $4, gensub(/(.*)\..*/,"\\1","g",$5) }'
#curl -X POST -H "Content-Type: application/json" -d '{"ips":["42.1.66.19"]}' 0.0.0.0:7001/update/v1/ips
#curl -X DELETE -H "Content-Type: application/json" -d '{"ips":["8.8.8.8","103.199.79.26","42.1.64.19","8.8.8.8"]}' 0.0.0.0:7001/delete/v1/ips
#htop -p $(pgrep -d',' -f "go run main.go")
#
#YHA:
#    "SrcUDP1": "10.51.46.101",
#    "SrcUDP2": "10.51.46.103",
