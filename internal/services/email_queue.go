package services

import "encoding/json"

const EmailQueueKey = "email_queue"

func EnqueueEmail(job EmailJob) error {
	payload, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return RedisClient.LPush(Ctx, EmailQueueKey, payload).Err()
}
