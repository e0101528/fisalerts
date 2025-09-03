#!/bin/bash
if [ "_$1" == "_" ]  ||  [ ! -e $1 ] ; then 
 echo
 echo Must supply json filename 
 echo
else 
 export INFLUX_TOKEN='OYpjptgIruw7uw2QNzQ1XfJ0IUawbPDSBrHA7IkjNnHyvEObmj1QNjLpf-UoxavJsUG3TWDxTLZ1JsdHePog2Q=='
 echo 
 curl -XPOST --header "Content-type: application/json" -d @${1} --header "Authorization: Token ${INFLUX_TOKEN}" http://localhost:8086/api/v2/checks
fi
