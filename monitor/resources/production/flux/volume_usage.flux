from(bucket: "telegraf") |> range (start: -15m) 
  |> filter (fn: (r) => r["_measurement"] == "disk" ) 
  |> filter (fn: (r) => r["_field"] == "used_percent" )
  |> last()
