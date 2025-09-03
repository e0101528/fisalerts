from(bucket: "_tasks")
  |> range(start: -60s )
  |> filter(fn: (r) => r["taskID"] == "0ef717d33c2cf000" or r["taskID"] == "0ef71628e9ecf000")
  |> filter(fn: (r) => r["_field"] == "logs")

