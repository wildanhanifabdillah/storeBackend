package services

import (
	"fmt"
	"time"
)

func GenerateOrderID() string {
	return fmt.Sprintf(
		"INV-%s-%d",
		time.Now().Format("20060102"),
		time.Now().UnixNano(),
	)
}
