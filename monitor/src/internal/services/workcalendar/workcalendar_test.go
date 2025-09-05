package workcalendar_test

import (
	"monitor/internal/services/workcalendar"
	"testing"
	"time"
)

func TestCtoD(t *testing.T) {

	tt, e := time.ParseInLocation("2006-01-02T15:04:05", "2006-01-02T15:38:05", time.Now().Local().Location())
	t.Logf("Error  %v", e)

	t.Logf("Time   %v", tt.Local().String())
	t.Logf("Hour   %v", tt.Local().Hour())
	t.Logf("Minute %v", tt.Local().Minute())

	dur := workcalendar.CToD(tt)
	tdur, _ := time.ParseDuration("15h38m")
	if dur != tdur {
		t.Errorf("time to duration conversion failed - %v != %v", dur, tdur)
	}
}
