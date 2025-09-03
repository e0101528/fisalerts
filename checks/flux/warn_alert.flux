import "influxdata/influxdb/monitor"
import "http"
import "json"
import "experimental"
import "date"
import "array"
import "join"

option task = {name: "eMail Alert", every: 1m, offset: 0s}

headers = {"Content-Type": "application/json"}
endpoint = http["endpoint"](url: "http://localhost:12345")
notification = {
    _notification_rule_id: "0efc3c6f0adb2000",
    _notification_rule_name: "eMail Alert",
    _notification_endpoint_id: "0ef7160e50c06000",
    _notification_endpoint_name: "Dummy Endpoint",
}

thistime = now()

//SOURCE FOR SPOOF OK DATA - INFLUXDB FAILS TO HANDLE EMPTY RIGHT TABLE on LEFT JOIN
spoofa = [
{_level: "ok",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984100000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:55:00.000Z), _type: "threshold", host: "fake" },
{_level: "ok",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984130000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:55:30.000Z), _type: "threshold", host: "fake" }
{_level: "warn",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984160000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:56:00.000Z), _type: "threshold", host: "fake" },
{_level: "warn",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984190000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:56:30.000Z), _type: "threshold", host: "fake" }
]

spoof = array.from(rows: spoofa)

msource = monitor["from"](start: -15m) 
 |> filter(fn: (r) => r["_check_id"] == "0ef717d339b1e000")                   //TEST CHECK
 |> filter(fn: (r) => r["host"] == "cibccaprodadm2.na1.prod.azure")           //SERVER FILTER

// SOURCE FROM MONITOR UNIONED WITH THE SPOOF OK DATA
source = union(tables: [spoof,msource])
|> group(columns: ["_measurement","_check_id","_check_name","_message","_level","_source_measurement","_start","_stop","_type","host"]) |> sort(columns: ["host","_time"])

//ALL THE CHECKS THAT ARE CURRENTLY IN WARN STATE
inwarn = source 
 |> group(columns: ["_measurement","_check_id","_check_name","_message","_source_measurement","_start","_stop","_type","host"])
 |> sort(columns: ["host","_time"])
 |> last(column: "_level")  
 |>  filter(fn: (r) =>  r["_level"] == "warn")
 |>  map (fn: (r) => ({
 r with
 _time: thistime
 })) |> group()

statuses = source
 |> keep(columns: ["_measurement","_check_id","_check_name","_message","_level","_source_measurement","_source_timestamp","_start","_stop","_type","_time","host"])
 |> map (fn: (r) => ({
 r with
 timestamp: r["_source_timestamp"] // /1000000000
 }))

//CREATE MULTIPLE STREAMS WITH EACH STATE AND ONE WITH ALL COMBINED. DROP _messages from ALL SINCE THE MESSAGE MAY DEPEND ON STATE.

all = statuses |> drop(columns: ["_message"])
 |> map (fn: (r) => ({
 r with
 _level: "any"
 })) |> sort(columns: ["host","_time"])

ok = statuses |>  filter(fn: (r) =>  r["_level"] == "ok")  |> sort(columns: ["host","_time"])
warn =  statuses |>  filter(fn: (r) =>  r["_level"] == "warn")  |> sort(columns: ["host","_time"])

odlast = ok   |> difference(columns: ["timestamp"] ) |> sort(columns: ["host","_time"])  |> last(column: "timestamp") |> group()
wdlast = warn |> difference(columns: ["timestamp"] )|> sort(columns: ["host","_time"])  |> last(column: "timestamp") |> group()

// THE BIG JOIN LAST ROW WHEN IN WARN STATE + WARN INTERVAL + OK INTERVAL 

iwlast = join.left(
              left: inwarn, 
              right: wdlast, 
  on: (l, r) => 
                l.host == r.host 
            and l._check_id == r._check_id
            and l._measurement == r._measurement
            and l._check_name == r._check_name
            and l._source_measurement == r._source_measurement
              ,
  as: (l, r) => ({
                _measurement: l._measurement, 
                host: l.host, 
                _check_id: l._check_id, 
                _check_name: l._check_name, 
                _source_timestamp: l._source_timestamp, 
                _source_measurement: l._source_measurement, 
                w_source_timestamp: r._source_timestamp, 
                warn: r.timestamp, 
                _start: l._start, 
                _stop: l._stop, 
                _time: l._time,
                _type: l._type, 
                _message: l._message, 
                _level: l._level
               }), ) |> group() 

last = join.left(
              left: iwlast, 
              right: odlast, 
  on: (l, r) => 
                l.host == r.host 
            and l._check_id == r._check_id
            and l._measurement == r._measurement
            and l._check_name == r._check_name
            and l._source_measurement == r._source_measurement
              ,
  as: (l, r) => ({
                _measurement: l._measurement, 
                host: l.host, 
                _check_id: l._check_id, 
                _check_name: l._check_name, 
                l_source_timestamp: l._source_timestamp, 
                _source_measurement: l._source_measurement, 
                w_source_timestamp: l.w_source_timestamp, 
                o_source_timestamp: r._source_timestamp, 
                warn: l.warn, 
                ok: r.timestamp, 
                _start: l._start, 
                _stop: l._stop, 
                _time: l._time,
                _type: l._type, 
                _message: l._message, 
                _level: l._level
               }), ) 
 |> filter(fn: (r) => r["host"] != "fake")           //remove spoof data after join()

alerts = last  |> filter(fn: (r) =>  not exists r["warn"]
                                     or
                                     ( exists r["ok"] and r["w_source_timestamp"] -  r["o_source_timestamp"] <  r["warn"] )
                        )
// jj = union(tables: [inwarn,wdlast,odlast])
notify = alerts |> range(start: -1m)
oclean = ok |> filter(fn: (r) => r["host"] != "fake")
wclean = warn  |> filter(fn: (r) => r["host"] != "fake")
odclean = odlast |> filter(fn: (r) => r["host"] != "fake")
wdclean = wdlast |> filter(fn: (r) => r["host"] != "fake")
iwclean = iwlast |> filter(fn: (r) => r["host"] != "fake")
