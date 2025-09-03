from(bucket: "telegraf")
  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
  |> filter(fn: (r) => r["_field"] == "StuckThreadCount")
  |> filter(fn: (r) => r["_measurement"] == "weblogic.ThreadPoolRuntime")
  |> group(columns: ["_field", "host", "_measurement"])
  |> aggregateWindow(every: v.windowPeriod, fn: max, createEmpty: false)

