import "date"
import "timezone"
option location = timezone.location(name: "America/Toronto")
from(bucket: "telegraf")
  |> range(start: -2m)
  |> filter(fn: (r) => r["_measurement"] == "haproxy")
  |> filter(fn: (r) => r["_field"] == "downtime_rate")
  |> filter(fn: (r) => r["proxy"] !~ /beuscm/)
  |> filter(fn: (r) => date.hour(t: r._time) > 21 or  date.hour(t: r._time) < 20)
  |> max()
