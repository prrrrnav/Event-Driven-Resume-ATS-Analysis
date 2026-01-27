package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	"resume-ats-platform/worker-ats/processor"
)

type ResumeParsedEvent struct{
	JobID string `json:"job_id"`
	ResumeText string `json:"resume_text"`
	JobDesc string `json:"job_description"`

}

type KafkaConsumer struct {
	reader *kafka.Reader
	processor *processor.ATSProcessor
}

func NewKafkaConsumer() *KafkaConsumer{
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string {"localhost:9092"},
		Topic: "resume.parsed",
		GroupID: "worker-ats-group",
	})

	return &KafkaConsumer{
		reader: reader,
		processor: processor.NewATSProcessor(),
	}
}

func (kc *KafkaConsumer) Start(ctx context.Context){
	log.Println("Ats Worker listening to resume.parsed")

	for{
		msg,err := kc.reader.ReadMessage(ctx)
		if err!= nil {
			log.Printf("Kafka read stopped: %v",err)
			return
		}

		go kc.handleMessage(msg)
	}
}

func (kc *KafkaConsumer) handleMessage(msg kafka.Message) {
	var event ResumeParsedEvent

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("Failed to parse message: %v", err)
		return
	}

	// 👇 THIS IS THE IMPORTANT LINE YOU WERE ASKING ABOUT
	kc.processor.Process(
		event.JobID,
		event.ResumeText,
		event.JobDesc,
	)
}
