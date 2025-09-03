#!/bin/bash
for fname in check/*.check ; do
 echo Processing $fname
 check=$(basename $fname)
 flux=$(echo $check | sed 's/\.check/\.flux/')
 json=$(echo $check | sed 's/\.check/\.json/')
 if [ -e "flux/$flux" ] ; then 
  query=$(cat flux/$flux | jq -Rsa | sed 's/\\/\\\\/g' |sed 's~\/~\\\/~g')
  cp  json/$json bak/$json.bak
  cat check/$check | sed "s/##QUERY##/${query}/" > json/$json
 else 
  echo Missing flux "flux/$flux"
 fi
done
