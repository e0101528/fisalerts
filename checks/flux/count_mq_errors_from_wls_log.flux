import "strings"

from(bucket: "telegraf")
  |> range(start: -12h)
  |> filter(fn: (r) => r["_measurement"] == "AdminLog")
  |> filter(fn: (r) => r["_field"] == "notice")
  |> filter(fn: (r)  => r["level"] == "ERROR")
  |> filter(fn: (r) => strings.containsStr(v: r._value, substr: "fismq") == true)
  |> count()
