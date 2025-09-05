import "strings"

from(bucket: "telegraf")
  |> range(start: -2m)
  |> filter(fn: (r) => r["host"] == "cibccaprodwls1.na1.prod.azure" or r["host"] == "cibccaprodwls2.na1.prod.azure" or r["host"] == "cibccaprodwls3.na1.prod.azure" or r["host"] == "cibccaprodwls4.na1.prod.azure")
  |> filter(fn: (r) => r["_measurement"] == "weblogic.JMSRuntime")
  |> filter(fn: (r) => r["_field"] == "MessagesCurrentCount" )
  |> filter(fn: (r) => strings.containsStr(v: r.wls_Name, substr: "DeadQueue") == false )
  |> filter(fn: (r) => strings.containsStr(v: r.wls_Name, substr: "FisFundsBoslOrderQueue") == false )
  |> filter(fn: (r) => strings.containsStr(v: r.wls_Name, substr: "FisFundsBoslRealtimeOrderQueue") == false )
  |> max()
