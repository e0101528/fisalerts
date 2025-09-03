import "influxdata/influxdb/monitor"
import "http"
import "json"
import "experimental"
from(bucket: "_monitoring")
  |> range(start: -6m)
  |> drop(columns: ["cpu","_type","_check_name"] )
  |> filter(fn: (r) => r["_measurement"] == "statuses")
  |> filter(fn: (r) => r["host"] == "cibccaprodha2.na1.prod.azure")
  |> filter(fn: (r) => r["_level"] == "warn")
  |> filter(fn: (r) => r["_field"] == "_message")
//|> yield(name: "value")

from(bucket: "_monitoring")
  |> range(start: -6m)
  |> drop(columns: ["cpu","_type","_check_name"] )
  |> filter(fn: (r) => r["_measurement"] == "statuses")
  |> filter(fn: (r) => r["host"] == "cibccaprodha2.na1.prod.azure")
  |> filter(fn: (r) => r["_level"] == "warn")
  |> filter(fn: (r) => r["_field"] == "_message")
  |> group(columns: ["_field","_measurement","_source_measurement","host"])
  |> stateCount(fn: (r) => r._level == "warn", column: "age")
  // |> filter(fn: (r) => r["age"] > 3)
  |> max(column: "age")
  //    |> group()
  //|> filter(fn: (r) => r["age"] <= 4)
  //|> filter(fn: (r) => r["_time"] >= experimental["subDuration"](from: now(), d: 60s))
|> yield(name: "final")
