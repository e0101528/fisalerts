from(bucket: "telegraf")
  |> range(start: -5m)
  |> filter(fn: (r) => r["_measurement"] == "net_response")
  |> filter(fn: (r) => r["service"] =~ /Feedhub/)
  |> filter(fn: (r) => r["_value"] != "success")
  |> filter(fn: (r) => r["_field"] == "result")
  |> count()

