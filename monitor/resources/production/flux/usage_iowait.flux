from(bucket: "telegraf") 
  |> range (start: -3m) 
  |> filter (fn: (r) => r["_measurement"] == "cpu") 
  |> filter (fn: (r) =>  r["_field"] == "usage_iowait")
  |> min()
  
