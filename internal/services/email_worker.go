package services

import (
	"encoding/json"
	"log"
	"time"
)

func StartEmailWorker() {
	go func() {
		log.Println("📨 Email worker started")

		for {
			result, err := RedisClient.BLPop(Ctx, 0, EmailQueueKey).Result()
			if err != nil {
				log.Println("redis error:", err)
				time.Sleep(2 * time.Second)
				continue
			}

			var job EmailJob
			if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
				log.Println("invalid job:", err)
				continue
			}

			if err := SendPaymentSuccessEmailWithInvoice(
				job.To,
				job.OrderID,
				job.Amount,
				job.InvoicePath,
			); err != nil {
				log.Println("email send failed:", err)
			}
		}
	}()
}
