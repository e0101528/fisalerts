import "strings"

from(bucket: "telegraf")
  |> range(start: -2m)
  |> filter(fn: (r) => r["_measurement"] == "AdminLog")
  |> filter(fn: (r) => r["_field"] == "notice")
  |> filter(fn: (r)  => r["level"] == "ERROR")
  |> filter(fn: (r) => strings.containsStr(v: r._value, substr: "No job control record found") == true)
  |> duplicate(column: "_stop", as: "_time")
  |> count()
