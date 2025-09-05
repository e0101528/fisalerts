package repositories

import (
	"context"
	"fmt"
	"monitor/internal/utils"
	"strings"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Influxdb2 struct {
	Client influxdb2.Client
	Org    string
}

func InitInflux2(url string, auth string, org string) Influxdb2 {
	client := influxdb2.NewClient(url, auth)

	return Influxdb2{Client: client, Org: org}
}

func (I *Influxdb2) RunQuery(ctx context.Context, f string, t []string) map[string]interface{} {
	var r map[string]interface{}
	r = make(map[string]interface{})
	// Get query client
	queryAPI := I.Client.QueryAPI(I.Org)
	res, err := queryAPI.Query(ctx, f)
	if err == nil {
		for res.Next() {
			rec := res.Record()
			r[tagstokey(rec.Values(), t)] = rec.Values()
		}
	} else {
		utils.Error("RunQuery error: %v", err)
	}
	return r
}

func (I *Influxdb2) RunQueryWithParams(ctx context.Context, f string, p any) {

	queryAPI := I.Client.QueryAPI(I.Org)
	queryAPI.QueryWithParams(ctx, f, p)

}

func tagstokey(r map[string]interface{}, t []string) string {
	var keys []string
	for _, tag := range t {
		keys = append(keys, fmt.Sprintf("%v", r[tag]))
	}
	return strings.Join(keys, "_")
}
