from(bucket: "telegraf")
  |> range(start: -5m)
  |> filter(fn: (r) => r["_measurement"] == "dummy")
  |> filter(fn: (r) => r["_field"] == "dummyerror")
  |> max()


