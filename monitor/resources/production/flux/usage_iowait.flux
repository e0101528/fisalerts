from(bucket: "telegraf") 
  |> range (start: -5m) 
  |> filter (fn: (r) => r["_measurement"] == "cpu") 
  |> filter (fn: (r) =>  r["_field"] == "usage_iowait")
  |> min()
  
