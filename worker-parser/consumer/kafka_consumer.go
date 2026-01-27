package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"resume-ats-platform/worker-parser/parser"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
	writer *kafka.Writer
}

// Message schema for resume.uploaded
type ResumeUploadedEvent struct {
	JobID       string `json:"job_id"`
	ResumeBytes []byte `json:"resume_bytes"`
	FileType    string `json:"file_type"` // pdf / docx
	JobDesc     string `json:"job_description"`
}

type ResumeParsedEvent struct {
	JobID       string `json:"job_id"`
	ResumeText  string `json:"resume_text"`
	JobDesc     string `json:"job_description"`
	ParsedAtUTC time.Time `json:"parsed_at_utc"`
}

func NewKafkaConsumer() (*KafkaConsumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "resume.uploaded",
		GroupID:"worker-parser-group",
	})

	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "resume.parsed",
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaConsumer{
		reader: reader,
		writer: writer,
	}, nil
}

func (kc *KafkaConsumer) Start(ctx context.Context) error {
	log.Println("Kafka consumer started (resume.uploaded)")

	for {
		msg, err := kc.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		go kc.handleMessage(ctx, msg)
	}
}

func (kc *KafkaConsumer) handleMessage(ctx context.Context, msg kafka.Message) {
	var event ResumeUploadedEvent

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	log.Printf("Processing resume for job_id=%s", event.JobID)

	// Call parser (pure logic)
	parsedText, err := parser.ParseResume(event.ResumeBytes, event.FileType)
	if err != nil {
		log.Printf("Resume parsing failed for job_id=%s: %v", event.JobID, err)
		return
	}

	outEvent := ResumeParsedEvent{
		JobID:       event.JobID,
		ResumeText:  parsedText,
		JobDesc:     event.JobDesc,
		ParsedAtUTC: time.Now().UTC(),
	}

	payload, err := json.Marshal(outEvent)
	if err != nil {
		log.Printf("Failed to marshal parsed event: %v", err)
		return
	}

	err = kc.writer.WriteMessages(ctx, kafka.Message{
		Value: payload,
	})
	if err != nil {
		log.Printf("Failed to publish resume.parsed event: %v", err)
		return
	}

	log.Printf("Published resume.parsed for job_id=%s", event.JobID)
}
