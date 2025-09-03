#!/bin/bash
if [ "_$1" == "_" ]  ; then 
 echo
 echo Must supply check ID
 ./list_checks.bash
 echo
else 
 export INFLUX_TOKEN='OYpjptgIruw7uw2QNzQ1XfJ0IUawbPDSBrHA7IkjNnHyvEObmj1QNjLpf-UoxavJsUG3TWDxTLZ1JsdHePog2Q=='
 echo 
 if [ "_$2" == "_full" ] ; then 
 curl --header "Authorization: Token ${INFLUX_TOKEN}" http://localhost:8086/api/v2/checks/${1} 2>/dev/null | jq 
else
 curl --header "Authorization: Token ${INFLUX_TOKEN}" http://localhost:8086/api/v2/checks/${1} 2>/dev/null | jq '.|{name, orgID, query:{text: .query.text, name: .query.name}, statusMessageTemplate, every, offset, tags, thresholds,type,status}' 
fi
fi
