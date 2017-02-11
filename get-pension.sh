#!/bin/bash

chain=$1
#echo $chain 
#echo curl -X POST --data-binary '{"request": "pension", "params":"'${chain}'"}' -H 'content-type:text/plain;' http://localhost:1337/POST
curl -X POST --data-binary '{"request": "pension", "params":"'${chain}'"}' -H 'content-type:text/plain;' http://localhost:1337/POST
# curl -X POST --data-binary '{"request": "all-pensions"}' -H 'content-type:text/plain;' http://localhost:1337/POST