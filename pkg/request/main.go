package request

import (
	"time"

	"github.com/go-resty/resty/v2"
)

var Client = resty.New()

func init() {
	Client.
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second)
}
