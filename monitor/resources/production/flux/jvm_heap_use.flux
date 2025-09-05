import "experimental"

from(bucket: "telegraf")
  |> range(start: -3m)
  |> filter(fn: (r) => r["host"] == "cibccaprodwls1.na1.prod.azure" or r["host"] == "cibccaprodwls2.na1.prod.azure" or r["host"] == "cibccaprodwls3.na1.prod.azure" or r["host"] == "cibccaprodwls4.na1.prod.azure")
  |> filter(fn: (r) => r["_measurement"] == "weblogic.jvm_memory" )
  |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
  |> map (
  fn: (r) => ({
 r with 
  HeapUse:  r["HeapMemoryUsage.used"] / r["HeapMemoryUsage.committed"],
 })
 )
|> drop(columns: [ "HeapMemoryUsage.used", "HeapMemoryUsage.committed"])
|> experimental.unpivot()
  |> filter(fn: (r) => r["_field"] == "HeapUse" )
|> mean()
|> fill(value: 0.0)
