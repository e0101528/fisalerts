import "influxdata/influxdb/monitor"
import "http"
import "json"
import "experimental"

statuses = monitor["from"](start: -900s)
statuses |> monitor.stateChanges(fromLevel: "ok", toLevel: "warn")
//  |> filter(fn: (r) => r["_time"] >= experimental["subDuration"](from: now(), d: 290s))
