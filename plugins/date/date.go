package date

import (
	"context"
	"time"

	"code.cloudfoundry.org/dontpanic/osreporter"
)

func Run(ctx context.Context) ([]byte, error) {
	return osreporter.WithTimeout(ctx, func() ([]byte, error) {
		d := time.Now().Format(time.UnixDate)
		return []byte(d), nil
	})
}
