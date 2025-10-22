from(bucket: "telegraf")
  |> range(start: -12h)
  |> filter(fn: (r) => r["_measurement"] == "Bareos")
  |> filter(fn: (r) => r["_field"] == "status")
  |> filter(fn: (r) => r["_value"] != "ERROR")
  |> group(columns: ["_measurement","_field","host"])
  |> count()

