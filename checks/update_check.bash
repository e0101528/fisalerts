#!/bin/bash
if [ "_$1" == "_" ]  ||  [ "_$2" == "_" ]  ||  [ ! -e $2 ] ; then 
 echo
 echo Must supply check ID and json filename 
 ./list_checks.bash
 ls ./json
 echo
else 
 export INFLUX_TOKEN='OYpjptgIruw7uw2QNzQ1XfJ0IUawbPDSBrHA7IkjNnHyvEObmj1QNjLpf-UoxavJsUG3TWDxTLZ1JsdHePog2Q=='
 echo 
 curl -XPUT --header "Content-type: application/json" -d @${2} --header "Authorization: Token ${INFLUX_TOKEN}" http://localhost:8086/api/v2/checks/${1}
fi
