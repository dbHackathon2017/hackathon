#!/bin/bash

chain=$1
#echo $chain 
#echo curl -X POST --data-binary '{"request": "pension", "params":"'${chain}'"}' -H 'content-type:text/plain;' http://localhost:1337/POST
curl -X POST --data-binary '{"request": "pension", "params":"'${chain}'"}' -H 'content-type:text/plain;' http://localhost:1337/POST
# curl -X POST --data-binary '{"request": "all-pensions"}' -H 'content-type:text/plain;' http://localhost:1337/POST
# curl -X POST --data-binary '{"request": "transaction", "params":"9b28ebf5a0ff896973b8ba780ebe43f52cde98b57278d9d5573bfcontent-type:text/plain;' http://localhost:1337/POST
