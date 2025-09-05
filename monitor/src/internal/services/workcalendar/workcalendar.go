package workcalendar

import (
	//"monitor/internal/config"

	"monitor/internal/utils"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rickar/cal/v2"
	"github.com/rickar/cal/v2/ca"
	"github.com/rickar/cal/v2/us"
)

var Weekdays = map[string]time.Weekday{
	"saturday":  time.Saturday,
	"sunday":    time.Sunday,
	"friday":    time.Friday,
	"thursday":  time.Thursday,
	"wednesday": time.Wednesday,
	"tuesday":   time.Tuesday,
	"monday":    time.Monday,
	"daily":     time.Weekday(8),
}

type Calendar struct {
	Name             string        `yaml:"name"`
	ActiveDays       []string      `yaml:"activedays"`
	ActiveHoursStart time.Duration `yaml:"activehoursstart"`
	ActiveHoursEnd   time.Duration `yaml:"activehoursend"`
	Holidays         []string      `yaml:"holidays"` // In rickar/cal/v2 Format
	Cal              *cal.BusinessCalendar
	MaintenanceDay   string        `yaml:"maintenanceday"`
	MaintenanceStart time.Duration `yaml:"maintenancestart"`
	MaintenanceStop  time.Duration `yaml:"maintenancestop"`
}

type CalendarMap struct {
	c map[string]*Calendar
}

func InitCalendarService(cfgcals map[string]Calendar) *CalendarMap {
	s := ""
	c := make(map[string]*Calendar)
	c[s] = &Calendar{Name: "FIS Operations Default 24x7"}

	c[s].Cal = cal.NewBusinessCalendar()
	c[s].Cal.SetWorkday(time.Saturday, true)
	c[s].Cal.SetWorkday(time.Sunday, true)
	c[s].Cal.SetWorkHours(0*time.Hour, 23*time.Hour+59*time.Minute)
	c[s].Cal.AddHoliday(
		ca.NewYear,
		us.LaborDay,
	)
	for nm, cl := range cfgcals {
		c[nm] = &Calendar{
			Name:             cl.Name,
			MaintenanceDay:   cl.MaintenanceDay,
			MaintenanceStart: cl.MaintenanceStart,
			MaintenanceStop:  cl.MaintenanceStop,
		}
		utils.Info("Initializing Calendar %s [name: %s]\n", nm, c[nm].Name)

		c[nm].Cal = cal.NewBusinessCalendar()
		if len(cl.ActiveDays) > 0 {
			for wdn, d := range Weekdays {
				if wdn != "daily" {
					c[nm].Cal.SetWorkday(d, false)
				}
			}
			for _, s := range cl.ActiveDays {
				d, ok := Weekdays[strings.ToLower(s)]
				if ok {
					c[nm].Cal.SetWorkday(d, true)
				}
			}
		} else {
			for _, d := range Weekdays {
				c[nm].Cal.SetWorkday(d, true)
			}
		}
		c[nm].Cal.SetWorkday(time.Saturday, true)
		c[nm].Cal.SetWorkday(time.Sunday, true)
		c[nm].Cal.SetWorkHours(cl.ActiveHoursStart, cl.ActiveHoursEnd)
		for _, h := range cl.Holidays {
			holiday, err := GetHolidayByName(h)
			if err == nil {
				c[nm].Cal.AddHoliday(holiday)
			}
		}
	}
	return &CalendarMap{
		c: c,
	}
}

func (c *CalendarMap) IsActive(calendarName string) bool {
	cal, ok := c.c[calendarName]
	if ok && cal != nil {
		utils.Info("Calendar Name: %s", cal.Name)
		return cal.Cal.IsWorkTime(time.Now())
	} else {
		utils.Warn("Invalid Calendar Name")
		return true //we don't care if there is an error, just assume active
	}
}

func GetHolidayByName(name string) (*cal.Holiday, error) {
	for _, h := range ca.Holidays {
		if h.Name == name {
			return h, nil
		}

	}
	for _, h := range us.Holidays {
		if h.Name == name {
			return h, nil
		}
	}
	return nil, errors.New("Unknown Holiday")

}

func (c *CalendarMap) InMaintenance(calendarName string) bool {
	cal, ok := c.c[calendarName]
	if ok && cal != nil {
		utils.Info("InMaintenance Calendar Name: %s", cal.Name)
		d, ok := Weekdays[strings.ToLower(cal.MaintenanceDay)]
		if !ok {
			utils.Info("InMaintenance invalid maintenance day: %s in %s", cal.MaintenanceDay, cal.Name)

			return false
		}
		n := cToD(time.Now())
		if (d == 8 || time.Now().Weekday() == d) && cal.MaintenanceStart <= n && cal.MaintenanceStop >= n {

			return true
		} else {
			utils.Debug("InMaintenance: %d==%d && %s <= %s <= %s is false\n", time.Now().Weekday(), d, cal.MaintenanceStart.String(), n.String(), cal.MaintenanceStop.String())
		}
	} else {
		utils.Warn("InMaintenance Invalid Calendar Name")

	}
	utils.Info("InMaintenance not in maintenance")
	return false //we don't care if there is an error, just assume active
}

func cToD(t time.Time) (d time.Duration) {
	d = (time.Duration(t.Local().Hour()) * time.Hour) + (time.Duration(t.Local().Minute()) * time.Minute)
	return d
}
