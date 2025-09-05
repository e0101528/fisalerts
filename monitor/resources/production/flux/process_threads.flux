from(bucket: "telegraf") 
  |> range (start: -1m) 
  |> filter (fn: (r) => r["_measurement"] == "processes" ) 
  |> filter (fn: (r) => r["_field"] == "total_threads" ) 
  |> max()

