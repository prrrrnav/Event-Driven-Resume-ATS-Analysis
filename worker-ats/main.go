package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"resume-ats-platform/worker-ats/consumer"
	"syscall"
	"time"
)
func main(){
	log.Println("Starting ATS Worker Service")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kafkaConsumer := consumer.NewKafkaConsumer()

	go kafkaConsumer.Start(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received")

	cancel()
	time.Sleep(2 * time.Second)

	log.Println("ATS Worker stopped cleanly")
}