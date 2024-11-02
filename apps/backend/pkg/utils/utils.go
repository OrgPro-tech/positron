package utils

import (
	"fmt"
	"time"
)

func CreateOrderId() string {
	t := time.Now().UnixNano()
	return fmt.Sprintf("%v", t)
}
