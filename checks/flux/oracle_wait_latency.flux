from(bucket: "telegraf")
  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
  |> filter(fn: (r) => r["host"] == "cibccaproddb2.na1.prod.azure")
  |> filter(fn: (r) => r["instance"] == "PORA4")
  |> filter(fn: (r) =>
    r._measurement == "oracle_wait_event" and
    (r._field == "count" or r._field == "latency")
  )
  |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
  |> map(fn: (r) => ({
    r with
    _value: r.latency * r.count 
   }))
    |> group(columns: ["wait_class","instance","wait_event"])
    |> aggregateWindow(every: 1m, column: "_value", fn: max)
