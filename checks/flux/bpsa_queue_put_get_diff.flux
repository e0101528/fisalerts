import "join"
difference = (x, y) => x - y

left = from(bucket: "telegraf")
  |> range(start: v.timeRangeStart, stop:v.timeRangeStop)
  |> filter(fn: (r) =>
    r._measurement == "queue" and
    r._field == "mqget_count" and
    r.queue == "CIBC.BPSA.A066.RTXML" 
  )
        |> group(columns: ["_time", "_value", "_field"], mode: "except")

right = from(bucket: "telegraf")
  |> range(start: v.timeRangeStart, stop:v.timeRangeStop)
  |> filter(fn: (r) =>
    r._measurement == "queue" and
    r._field == "mqput_mqput1_count" and
    r.queue == "CIBC.BPSA.A066.RTXML" 
  )
          |> group(columns: ["_time", "_value", "_field"], mode: "except")


join.time(method: "left", left: left, right: right, as: (l, r) => ({l with target: r._value - l._value}))

|> filter(fn: (r) => 
r.target > 0
)
