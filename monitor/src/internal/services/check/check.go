package check

import (
	"context"
	"fmt"
	"monitor/internal/config"
	infl "monitor/internal/repositories/influxdb"
	"monitor/internal/utils"
	"time"

	"strings"

	"github.com/google/uuid"
)

var levels []string = []string{"", "Critical", "Warning", "Information", "OK"}

func Run(ctx context.Context, c config.Check) {

	var result bool
	utils.Info("Running check %s\n", c.Name)
	cfg := ctx.Value("config").(*config.ApplicationConfig)

	r := infl.InitInflux2(cfg.InfluxDB.ServerURL, cfg.InfluxDB.AuthToken, cfg.InfluxDB.Org)
	res := r.RunQuery(ctx, c.Flux, c.Tags)

	db := cfg.Storage
	cal := cfg.CalendarMap
	nfys := cfg.Notifiers
	utils.Dump("Check", c)
	nfy, ok := nfys[c.Target]
	if !ok {
		utils.Warn("Unknown notification target: %s\n", c.Target)
	}
	utils.Info("Got %d results\n", len(res))
	utils.Dump("Results", res)
	for k, v := range res {
		value := v.(map[string]interface{})["_value"]
		switch value.(type) {
		case string:
			result = compare(value, c.Comparison, c.Match)
		default:
			result = compare(value, c.Comparison, c.Threshold)

		}

		if result {
			switch value.(type) {
			case string:
				utils.Debug("Check %s is active:\n  string => \"%v\" %v \"%v\" \n", c.Name, value, c.Comparison, c.Match)
			default:
				utils.Debug("Check %s is active:\n  numeric => %v %v %v \n", c.Name, value, c.Comparison, c.Threshold)
			}
			if db.SetActive(ctx, []byte(c.Name), []byte(k), int64(c.Interval)) {
				if cal.IsActive(c.Calendar) {
					content := RtoC(cfg, v.(map[string]interface{}), c)
					utils.Debug("!!! AWWOOGA %s \n", k)
					utils.Dump("Value", value)
					if ok {
						if !cal.InMaintenance(c.Calendar) {
							nfy.SendAlert(content)
						} else {
							utils.Debug("In maintenance mode\n")
						}
					}
				} else {
					utils.Info("Calendar %s inactivates check %s at the moment [%s]\n", c.Calendar, c.Name, time.Now().String())

				}

			}
		}
	}
	//utils.Dumper(res)
}
func RtoC(cfg *config.ApplicationConfig, r map[string]interface{}, c config.Check) (content map[string]string) {

	content = make(map[string]string)
	for n, v := range r {
		content[n] = fmt.Sprintf("%v", v)

	}

	content["self"] = cfg.Webserver.Host
	content["status"] = "firing"
	content["severity"] = fmt.Sprintf("SEV-%d", c.Severity)
	content["_level"] = levels[c.Severity]

	//	content["self"] = o.Self
	content["uuid"] = uuid.NewString()
	content["resourceGroup"] = cfg.AzureRG
	content["subscription"] = cfg.AzureSub
	content["environment"] = "Test"
	content["appGroupEmail"] = cfg.AppGroupEmail
	content["assignmentGroup"] = cfg.AssignmentGroup
	content["ipaddr"] = cfg.IPAddr
	content["recipients"] = cfg.AssignmentGroup
	content["org"] = cfg.Organization
	content["_type"] = "Threshold"
	content["_check_name"] = c.Name
	content["_check_id"] = c.ID

	//utils.Dumper(content)
	return content
}

func compare(value interface{}, comparator string, match interface{}) bool {
	var num float64
	var str string
	var isnumneric bool

	switch v := value.(type) {
	case nil:
		utils.Debug("x is nil") // here v has type interface{}
	case int:
		num = float64(v)
		isnumneric = true
	case int8:
		num = float64(v)
		isnumneric = true
	case int16:
		num = float64(v)
		isnumneric = true
	case int32:
		num = float64(v)
		isnumneric = true
	case int64:
		num = float64(v)
		isnumneric = true
	case float32:
		num = float64(v)
		isnumneric = true
	case float64:
		num = v
		isnumneric = true
	case bool:
		if v {
			num = 1
		} else {
			num = 0
		}
		isnumneric = true
	case string:
		str = v
		utils.Debug("value is of type string\n") // here v has type interface{}
		isnumneric = false
	default:
		utils.Debug("value is of type unknown\n") // here v has type interface{}
		return false
	}
	if isnumneric {
		utils.Debug("numeric comparison\n") // here v has type interface{}

		threshold, ok := match.(float64)
		if ok {
			switch comparator {
			case "gt":
				return num > threshold
			case "lt":
				return num < threshold
			case "le":
				return num <= threshold
			case "ge":
				return num >= threshold
			case "ne":
				return num != threshold
			case "eq":
				return num == threshold
			default:
				utils.Warn("unknown numeric comparison %s\n", comparator)
				return false
			}
		}
	} else {
		matchstring, ok := match.(string)

		if ok {
			switch comparator {
			case "ne":
				return matchstring != str
			case "substr":
				utils.Debug("does \"%s\" contain \"%s\"\n", str, matchstring)
				return strings.Contains(str, matchstring)
			case "eq":
				return matchstring == str
			default:
				utils.Debug("unknown string comparison %s\n", comparator)
				return false
			}
		}
	}
	return false
}
