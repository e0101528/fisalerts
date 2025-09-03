import "influxdata/influxdb/monitor"
import "http"
import "json"
import "experimental"
import "date"
import "array"
import "join"

//TEST DATA
rows = [
{_level: "ok",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984160000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:56:00.000Z), _type: "threshold", host: "host1" },
{_level: "warn",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984190000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:56:30.000Z), _type: "threshold", host: "host1" },
{_level: "ok",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984220000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:57:00.000Z), _type: "threshold", host: "host1" },
{_level: "warn",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984250000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:57:30.000Z), _type: "threshold", host: "host1" },

{_level: "ok",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984190000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:56:30.000Z), _type: "threshold", host: "host2" },
{_level: "warn",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984220000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:57:00.000Z), _type: "threshold", host: "host2" },
{_level: "warn",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984250000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:57:30.000Z), _type: "threshold", host: "host2" },

{_level: "ok",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984190000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:56:30.000Z), _type: "threshold", host: "host3" },
{_level: "ok",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984220000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:57:00.000Z), _type: "threshold", host: "host3" },
{_level: "warn",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984250000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:57:30.000Z), _type: "threshold", host: "host3" },

{_level: "ok",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984190000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:56:30.000Z), _type: "threshold", host: "host4" },
{_level: "warn",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984220000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:57:00.000Z), _type: "threshold", host: "host4" },
{_level: "ok",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984250000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:57:30.000Z), _type: "threshold", host: "host4" },

{_level: "warn",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984220000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:57:00.000Z), _type: "threshold", host: "host5" },
{_level: "ok",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984250000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:57:30.000Z), _type: "threshold", host: "host5" },

{_level: "ok",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984160000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:56:00.000Z), _type: "threshold", host: "host6" },
{_level: "warn",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984190000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:56:30.000Z), _type: "threshold", host: "host6" }
]
// WARNING FOR host1 host3 host6
thistime = now()

//SOURCE FOR SPOOF OK DATA - INFLUXDB FAILS TO HANDLE EMPTY RIGHT TABLE on LEFT JOIN
 
spoofa = [
{_level: "ok",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984160000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:56:00.000Z), _type: "threshold", host: "fake" },
{_level: "ok",_measurement: "statuses",_check_id: "0ef717d339b1e000",_check_name: "Dummy",_message: "no msg",_source_measurement: "dummy",_source_timestamp: 1748984190000000000, _start: date.time(t: 2025-06-03T19:58:30.496Z),_stop: date.time(t: 2025-06-03T20:58:30.496Z),_time: date.time(t: 2025-06-03T20:56:30.000Z), _type: "threshold", host: "fake" }
]

spoof = array.from(rows: spoofa)


msource = monitor["from"](start: -10m) 
// |> filter(fn: (r) => r["_check_id"] == "0ef717d339b1e000")                   //TEST CHECK
// |> filter(fn: (r) => r["host"] == "cibccaprodadm2.na1.prod.azure")           //SERVER FILTER

// SOURCE FROM MONITOR UNIONED WITH THE SPOOF OK DATA
//source = union(tables: [spoof,msource])
//|> group(columns: ["_measurement","_check_id","_check_name","_message","_level","_source_measurement","_start","_stop","_type","host"]) |> sort(columns: ["host","_time"])

//SOURCE FOR TEST DATA
source = array.from(rows: rows) |> group(columns: ["_measurement","_check_id","_check_name","_message","_level","_source_measurement","_start","_stop","_type","host"]) |> sort(columns: ["host","_time"])

//ALL THE CHECKS THAT ARE CURRENTLY IN WARN STATE
inwarn = source 
 |> group(columns: ["_measurement","_check_id","_check_name","_message","_source_measurement","_start","_stop","_type","host"])
 |> sort(columns: ["host","_time"])
 |> last(column: "_level")  
 |>  filter(fn: (r) =>  r["_level"] == "warn")
 |>  map (fn: (r) => ({
 r with
 _time: thistime
 }))
 |> group()

//ALL THE CHECKS THAT ARE CURRENTLY IN CRIT STATE
incrit = source 
 |> group(columns: ["_measurement","_check_id","_check_name","_message","_source_measurement","_start","_stop","_type","host"])
 |> sort(columns: ["host","_time"])
 |> last(column: "_level")  
 |>  filter(fn: (r) =>  r["_level"] == "crit")
 |>  map (fn: (r) => ({
 r with
 _time: thistime
 }))
 |> group()


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
crit =  statuses |>  filter(fn: (r) =>  r["_level"] == "crit")  |> sort(columns: ["host","_time"])
info =  statuses |>  filter(fn: (r) =>  r["_level"] == "info")  |> sort(columns: ["host","_time"])

//CALCULATE INTERVAL BETWEEN LAST AND PREVIOUS STATE
adlast = all  |> difference(columns: ["timestamp"] ) |> sort(columns: ["host","_time"])  |> last(column: "timestamp") |> group()
odlast = ok   |> difference(columns: ["timestamp"] ) |> sort(columns: ["host","_time"])  |> last(column: "timestamp") |> group()
wdlast = warn |> difference(columns: ["timestamp"] )|> sort(columns: ["host","_time"])  |> last(column: "timestamp") |> group()
cdlast = crit |> difference(columns: ["timestamp"] )|> sort(columns: ["host","_time"])  |> last(column: "timestamp") |> group()

//UNUSED
//alast = all  |> last(column: "timestamp")
//olast = ok   |> last(column: "timestamp")
//wlast = warn |> last(column: "timestamp")
//clast = crit |> last(column: "timestamp")


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
               }), ) 

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

alerts = last  |> filter(fn: (r) =>  not exists r["warn"]
                                     or
                                     ( exists r["ok"] and r["w_source_timestamp"] -  r["o_source_timestamp"] <  r["warn"] )
                        )

jj = union(tables: [inwarn,wdlast,odlast])
// |> drop(columns: ["_start","_stop","_time"])

alerts

