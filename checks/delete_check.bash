#!/bin/bash
if [ "_$1" == "_" ] ; then 
 echo
 echo Must supply Job ID 
 echo ==================
 ./list_checks.bash
else 
 export INFLUX_TOKEN='OYpjptgIruw7uw2QNzQ1XfJ0IUawbPDSBrHA7IkjNnHyvEObmj1QNjLpf-UoxavJsUG3TWDxTLZ1JsdHePog2Q=='
 echo 
 curl -XDELETE --header "Authorization: Token ${INFLUX_TOKEN}" http://localhost:8086/api/v2/checks/$1
 echo Deleted $1
fi
