from(bucket: "telegraf")
  |> range(start: -2m)
  |> filter(fn: (r) => r["host"] == "cibccaprodwls1.na1.prod.azure" or r["host"] == "cibccaprodwls2.na1.prod.azure" or r["host"] == "cibccaprodwls3.na1.prod.azure" or r["host"] == "cibccaprodwls4.na1.prod.azure")
  |> filter(fn: (r) => r["_measurement"] == "weblogic.ServerRuntime")
  |> filter(fn: (r) => r["_field"] == "OpenSocketsCurrentCount" )
  |> mean()
  
