import "strings"

from(bucket: "telegraf")
  |> range(start: -5m)
  |> filter(fn: (r) => r["_measurement"] == "AdminLog")
  |> filter(fn: (r) => r["_field"] == "notice")
  |> filter(fn: (r)  => r["level"] == "ERROR")
  |> filter(fn: (r) => strings.containsStr(v: r._value, substr: "Cannot submit job as allow concurrent is set to N for job CFG.USCMUserSync") == true)
  |> count()
