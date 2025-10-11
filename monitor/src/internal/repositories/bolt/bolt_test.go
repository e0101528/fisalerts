package repositories

import (
	"context"
	"monitor/internal/utils"
	"testing"
	"time"
)

func TestInitBolt(t *testing.T) {
	ts, e := InitBolt()
	if e != nil {
		t.Error("Failed to initialize bolt DB")
	} else {
		ts.Shutdown()
	}

}

func TestSetActive(t *testing.T) {
	utils.SetLevel(6)

	ts, e := InitBolt() //debug
	if e != nil {
		t.Error("Failed to initialize bolt DB")
	} else {
		b1 := ts.SetActive(context.Background(), []byte("test"), []byte("testkey"), 5.0, map[string]any{"foo": "bar", "bar": 1.0})
		time.Sleep(2 * time.Second)
		b2 := ts.SetActive(context.Background(), []byte("test"), []byte("testkey"), 5.0, map[string]any{"foo": "bar", "bar": 1.0})
		time.Sleep(7 * time.Second)
		b3 := ts.SetActive(context.Background(), []byte("test"), []byte("testkey"), 5.0, map[string]any{"foo": "bar", "bar": 1.0})
		if !b1 {
			t.Error("initial state was active")
		}
		if b2 {
			t.Error("initial state was cleared too soon")
		}
		if !b3 {
			t.Error("second state was cleared too late")
		}
		ts.Shutdown()
	}

}
