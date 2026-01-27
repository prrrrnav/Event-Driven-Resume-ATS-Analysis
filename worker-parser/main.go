package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"resume-ats-platform/worker-parser/consumer"
)

func main(){
	log.Println("Starting Worker Parser Service")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kafkaConsumer, err := consumer.NewKafkaConsumer()
	if err!= nil{
		log.Fatalf("Failed to create kafka consumer : %v", err)
	}

	go func(){
		if err := kafkaConsumer.Start(ctx); err!= nil{
			log.Fatalf("Kafka consumer stopped with error:%v",err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received")
	cancel()
	time.Sleep(2* time.Second)

	log.Println("Worker parser Service stopped cleanly")
}