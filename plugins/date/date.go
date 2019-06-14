package date

import (
	"time"
)

func Run() ([]byte, error) {
	d := time.Now().Format(time.UnixDate)
	return []byte(d), nil
}
