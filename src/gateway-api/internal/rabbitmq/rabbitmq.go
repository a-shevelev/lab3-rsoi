package rabbitmq

import (
	"encoding/json"
	"errors"
	"fmt"
	"gateway-api/internal/dto"
	"gateway-api/pkg/ext"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func RunRetryWorker(ch *amqp.Channel, queueName string, process func(evt dto.ReturnRetryEvent) error) error {
	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to register consumer for %s: %w", queueName, err)
	}

	go func() {
		for d := range msgs {
			var evt dto.ReturnRetryEvent
			if err := json.Unmarshal(d.Body, &evt); err != nil {
				log.Printf("[worker %s] bad payload, acking: %v", queueName, err)
				_ = d.Ack(false)
				continue
			}

			for {
				err := process(evt)
				if err == nil {
					_ = d.Ack(false)
					log.Printf("[worker %s] success: %+v", queueName, evt)
					break
				}

				if errors.Is(err, ext.ServiceUnavailableError) {
					log.Printf("[worker %s] service unavailable, retrying in 10s: %v", queueName, err)
					time.Sleep(10 * time.Second)
					continue
				}

				log.Printf("[worker %s] non-retryable, dropping: %v", queueName, err)
				_ = d.Ack(false)
				break
			}
		}
	}()

	return nil
}
