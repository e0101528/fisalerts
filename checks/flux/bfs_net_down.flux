import "date"
import "timezone"
option location = timezone.location(name: "America/Toronto")

from(bucket: "telegraf")
  |> range(start: -5d)
  |> filter(fn: (r) => r["_measurement"] == "net_response")
  |> filter(fn: (r) => r["server"] == "149.83.1.11")
  |> filter(fn: (r) => r["port"] == "1364")
  |> filter(fn: (r) => r["result"] == "connection_failed")
  |> filter(fn: (r) => r["_field"] == "result_code")
  |> filter(fn: (r) => date.hour(t: r._time) > 21 or  date.hour(t: r._time) < 20)
  |> aggregateWindow(every: 5m, fn: count, createEmpty: false)
