#!/bin/bash
dst=$(date -u +%FT%TZ)
sst=$(date -u +%s%N)
sleep 1
dsp=$(date -u +%FT%TZ)
ssp=$(date -u +%s%N)
fo=$(mktemp)
cat >$fo <<EOJ
{
  "_check_id": "abc123defabc1234",
  "_check_name": "Influx relay to xMatters Test",
  "_level": "warn",
  "_measurement": "notifications",
  "_message": "Acknowledge and ignore. This is a test message",
  "_notification_endpoint_id": "0dfeb75956902000",
  "_notification_endpoint_name": "NineOneOne",
  "_notification_rule_id": "0dfed9329c138123",
  "_notification_rule_name": "xMatters Incident",
  "_source_measurement": "Nothing",
  "_source_timestamp": ${sst},
  "_start": "${dst}",
  "_status_timestamp": ${ssp},
  "_stop": "${dsp}",
  "_time": "${dst}",
  "_type": "threshold",
  "_version": 1
}

EOJ
echo $fo
# curl -XPOST -d @$fo localhost:12345/alert
 curl -XPOST -d @$fo localhost:9911/alert
rm -f $fo
