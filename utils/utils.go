package utils

import "time"

// 备用
func Retry(fn func() error, retries int) error {
	var err error
	for i := 0; i < retries; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return err
}
