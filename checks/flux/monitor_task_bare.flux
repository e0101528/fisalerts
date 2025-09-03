import "influxdata/influxdb/monitor"
import "http"
import "json"
import "experimental"
from(bucket: "_monitoring")
  |> range(start: -2m)
//  |> drop(columns: ["cpu","_type","_check_name"] )
//  |> filter(fn: (r) => r["_measurement"] == "statuses")
  |> filter(fn: (r) => r["host"] == "cibccaprodha1.na1.prod.azure" or r["host"] == "cibccaprodha2.na1.prod.azure")
//  |> filter(fn: (r) => r["_level"] == "warn")
//  |> filter(fn: (r) => r["_field"] == "_message")
//|> sort(columns: ["host","_time"])
|> count()
