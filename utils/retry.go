package utils

import (
	"time"

	"gopkg.in/src-d/go-log.v1"
)

const (
	retries  = 10
	delay    = 10 * time.Millisecond
	truncate = 10 * time.Second
)

func Retry(f func() error) error {
	d := delay
	var i uint

	for ; ; i++ {
		err := f()
		if err == nil {
			return nil
		}

		if i == retries {
			return err
		}

		log.Errorf(err, "retrying in %v", d)
		time.Sleep(d)

		d = d * (1<<i + 1)
		if d > truncate {
			d = truncate
		}
	}
}
