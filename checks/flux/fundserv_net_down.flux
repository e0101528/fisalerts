import "date"
import "timezone"
option location = timezone.location(name: "America/Toronto")

from(bucket: "telegraf")
  |> range(start: -5m)
  |> filter(fn: (r) => r["_measurement"] == "net_response")
  |> filter(fn: (r) => r["server"] == "142.148.10.116")
  |> filter(fn: (r) => r["port"] == "1501")
  |> filter(fn: (r) => r["result"] == "connection_failed")
  |> filter(fn: (r) => r["_field"] == "result_code")
  |> filter(fn: (r) => date.hour(t: r._time) > 1 and date.hour(t: r._time) < 20)
  |> aggregateWindow(every: 5m, fn: count, createEmpty: false)

