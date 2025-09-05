import "array"
import "join"
import "math"

yesterday = from(bucket: "telegraf") |> range (start: -2d, stop: -1d) 
  |> filter (fn: (r) => r["_measurement"] == "disk" ) 
  |> filter (fn: (r) => r["_field"] == "used_percent" ) 

today   =   from(bucket: "telegraf") |> range (start: -1d)
  |> filter (fn: (r) => r["_measurement"] == "disk" ) 
  |> filter (fn: (r) => r["_field"] == "used_percent" ) 


tf = today |> first()
yf = yesterday |> first()
tl = today |> last()
yl = yesterday |> last()
t = union(tables: [tf, tl])
y = union(tables: [yf, yl])

yd = y |> difference() |> drop(columns: ["_start","_stop","_time","mode","fstype","_field","_measurement","device"])
td = t |> difference() |> drop(columns: ["_start","_stop","_time","mode","fstype","_field","_measurement","device"])
// Table: keys: [_field, _measurement, device, fstype, host, mode, path]

d = join.left(
              left: yd, 
              right: td,
  on: (l, r) => 
                l.host == r.host 
            and l.path == r.path
              ,
  as: (l, r) => ({
                host: l.host,
                path: l.path,
                yd: l._value,
                td: r._value
               }), ) 

d |> map (fn: (r) => ({
 r with
 ratio: if r.td > 0 and r.yd != 0 then 
    math.abs(x: r.td / r.yd) 
   else
    0.0
 })) 
  |> filter(fn: (r) => (r.ratio > 2))


