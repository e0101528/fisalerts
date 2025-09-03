export INFLUX_TOKEN='OYpjptgIruw7uw2QNzQ1XfJ0IUawbPDSBrHA7IkjNnHyvEObmj1QNjLpf-UoxavJsUG3TWDxTLZ1JsdHePog2Q=='
echo 
echo -e "[Job Name]\t\t\tID"
curl --header "Authorization: Token ${INFLUX_TOKEN}" http://localhost:8086/api/v2/checks 2>/dev/null  | jq -r '.checks[] | "[\(.name)]\t\t\(.id)"'
echo
