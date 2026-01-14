package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(broker string , topic string) *Producer {
	w:= &kafka.Writer{
		Addr: kafka.TCP(broker),
		Topic: topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{
		writer: w,
	}
}

func (p *Producer) Publish(ctx context.Context, event interface{}) error{
	bytes, err:= json.Marshal(event)
	if err!= nil {
		return err
	}
	log.Println("Publishing event to Kafka...")
	
	msg:= kafka.Message{
		Value: bytes,
	}



	if err:= p.writer.WriteMessages(ctx,msg); err != nil{
		log.Printf("Kafka publish failed: %v\n", err)
		return err
	}

	log.Println("Kafka event published successfully")
	return nil

}