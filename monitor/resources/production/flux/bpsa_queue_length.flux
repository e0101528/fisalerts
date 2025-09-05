from(bucket: "telegraf")
  |> range(start: -2m)
  |> filter(fn: (r) =>
    r._measurement == "queue" and
    r._field == "queue_depth" and
    r.queue == "CIBC.BPSA.A066.RTXML" 
  )
    |> max()
