package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"monitor/internal/utils"
	"time"

	"github.com/boltdb/bolt"
)

type TimeSeries struct {
	TSDB *bolt.DB
}

func InitBolt() (*TimeSeries, error) {

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists([]byte("states"))
		return e
	})
	if err != nil {
		return nil, err
	}
	return &TimeSeries{
		TSDB: db,
	}, nil
}

func (tsdb *TimeSeries) Shutdown() {
	tsdb.TSDB.Close()
}

func (tsdb *TimeSeries) SetActive(ctx context.Context, check []byte, key []byte, interval int64, rm map[string]any) bool {
	var lm map[string]any
	isActivated := false
	pIsActivated := &isActivated
	now := time.Now().Unix()
	//buf := new(bytes.Buffer)
	rm["timestamp"] = int64(now)
	//err := binary.Write(buf, binary.BigEndian, rm)
	bts, err := json.Marshal(rm)
	if err != nil {
		utils.Error("binary.Write failed: %v", err)
	}
	err = tsdb.TSDB.Update(func(tx *bolt.Tx) error {
		bucket, e := tx.CreateBucketIfNotExists(check)
		if e != nil {
			*pIsActivated = false
			return e
		}

		v := bucket.Get([]byte(key))
		var lastactive int64
		var ok bool
		if v != nil {
			e := json.Unmarshal(v, &lm)
			if e != nil {
				lastactive = 0
				utils.Error("binary unmarshal failed: %v", err)

			} else {
				lastactive, ok = lm["timestamp"].(int64)
				if !ok {
					lastactive = 0
				}
			}
		} else {
			lastactive = 0
		}
		utils.Info("Active? lastactive[%d] < now[%d] - interval[%d](+20%%)\n", lastactive, now, interval)
		if lastactive < now-int64(float64(interval)*1.2) { //safety margin against slow running query
			*pIsActivated = true
			utils.Info("IsActivated=true\n")
		} else {
			utils.Info("IsActivated=false\n")

		}

		bucket.Put([]byte(key), bts)
		return nil
	})
	if err != nil {
		utils.Error("bolt Setactive error: %v", err)
	}
	return isActivated
}

func (tsdb *TimeSeries) Fetch(ctx context.Context, check []byte, key string) (data map[string]any, err error) {
	data = make(map[string]any)

	tx, e := tsdb.TSDB.Begin(false)
	if e != nil {
		return data, e
	}
	defer tx.Rollback()
	bucket := tx.Bucket(check)
	if bucket == nil {
		err = fmt.Errorf("bucket %v not found", check)
		return
	}
	v := bucket.Get([]byte(key))
	if v != nil {
		err = json.Unmarshal(v, &data)
		if err != nil {
			utils.Error("binary unmarshal failed: %v", err)
			return
		}

	}

	return data, nil

}
