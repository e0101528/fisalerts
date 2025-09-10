package check

import (
	"context"
	"fmt"
	"monitor/internal/config"
	infl "monitor/internal/repositories/influxdb"
	"monitor/internal/services/notifier"
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
						if cal.InMaintenance(c.Calendar) || cfg.MaintenanceMode {
							utils.Debug("In maintenance mode - calendar[%v] - command line [%v]\n", cal.InMaintenance(c.Calendar), cfg.MaintenanceMode)
						} else {
							nfy.SendAlert(content)
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
func RtoC(cfg *config.ApplicationConfig, r map[string]interface{}, c config.Check) (content notifier.Content) {

	content.Labels = make(map[string]string)
	content.Fields = make(map[string]string)
	for n, v := range r {
		if isin(n, &c.Tags) {
			content.Labels[n] = fmt.Sprintf("%v", v)
		}
	}

	content.Fields["self"] = cfg.Webserver.Host
	content.Fields["status"] = "firing"
	content.Fields["severity"] = fmt.Sprintf("SEV-%d", c.Severity)
	content.Fields["_level"] = levels[c.Severity]

	//	content.Fields["self"] = o.Self
	content.Fields["uuid"] = uuid.NewString()
	content.Fields["resourceGroup"] = cfg.AzureRG
	content.Fields["subscription"] = cfg.AzureSub
	content.Fields["environment"] = "Test"
	content.Fields["appGroupEmail"] = cfg.AppGroupEmail
	content.Fields["assignmentGroup"] = cfg.AssignmentGroup
	content.Fields["ipaddr"] = cfg.IPAddr
	content.Fields["recipients"] = cfg.AssignmentGroup
	content.Fields["org"] = cfg.Organization
	content.Fields["_type"] = "Threshold"
	content.Fields["_check_name"] = c.Name
	content.Fields["_check_id"] = c.ID
	var ok bool
	content.Fields["host"], ok = content.Labels["host"]
	if !ok {
		content.Fields["host"] = cfg.Webserver.Host
	}
	content.Fields["_time"], ok = content.Labels["_time"]
	if !ok {
		content.Fields["_time"] = time.Now().String()
	}
	content.Fields["_start"], ok = content.Labels["_start"]
	if !ok {
		content.Fields["_start"] = time.Now().String()
	}
	content.Fields["_stop"], ok = content.Labels["_stop"]
	if !ok {
		content.Fields["_stop"] = time.Now().String()
	}

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

func isin(s string, a *[]string) bool {
	for i := range *a {
		if (*a)[i] == s {
			return true
		}
	}
	return false
}
