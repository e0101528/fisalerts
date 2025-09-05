package repositories

import (
	"context"
	"encoding/binary"
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

func (tsdb *TimeSeries) SetActive(ctx context.Context, check []byte, key []byte, interval int64) bool {
	isActivated := false
	pIsActivated := &isActivated
	now := time.Now().Unix()
	n := make([]byte, 8)
	binary.BigEndian.PutUint64(n, uint64(now))
	err := tsdb.TSDB.Update(func(tx *bolt.Tx) error {
		bucket, e := tx.CreateBucketIfNotExists(check)
		if e != nil {
			*pIsActivated = false
			return e
		}
		v := bucket.Get([]byte(key))
		var lastactive int64
		if v != nil {
			lastactive = int64(binary.BigEndian.Uint64(v))
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
		bucket.Put([]byte(key), n)
		return nil
	})
	if err != nil {
		utils.Error("bolt Setactive error: %v", err)
	}
	return isActivated
}
