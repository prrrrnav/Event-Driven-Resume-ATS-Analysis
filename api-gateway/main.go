package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"net/http"

	"resume-ats-platform/api-gateway/handler"
	"resume-ats-platform/api-gateway/kafka"
)


func main(){
	producer := kafka.NewProducer(
		"localhost:9092",
		"resume.uploaded",
	)

	submitHandler := &handler.SubmitHandler{
		Producer:producer,
	}

	http.Handle("/submit",submitHandler)

	server:= &http.Server{
		Addr: ":8080",
	}

	go func(){
		log.Println("API Gateway running on :8080")
		if err := server.ListenAndServe(); err!= nil && err != http.ErrServerClosed{
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	
	<-stop
	log.Println(("Shutting down api gateway.."))

	ctx, cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	if err:= server.Shutdown(ctx);err != nil{
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("Server stopped cleanly")





}