package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	Id string
	number1 int32
	message string
}

func main() {
	cfg := new(tls.Config)
	cfg.RootCAs = x509.NewCertPool()
	conn, err := amqp.DialTLS("amqps://guest:guest@rabbitmq:5672/", cfg)
	if err != nil {
		log.Printf("tls no connection: %v\n", err)
	}

	conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/test")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare("test", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = ch.QueueBind("test", "test-routing-key", "amq.topic", false, nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(Message{
		Id: "1",
		number1: 42,
		message: "A simple message",
	})
	if err != nil {
		log.Fatal(err)
	}

	for _ = range make([]int, 5, 5) {
		err = ch.PublishWithContext(ctx,
			"amq.topic",
			"test-routing-key",
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				Timestamp: time.Now(),
				ContentType: "application/json",
				Body: body,
			},
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = ch.Qos(2, 0, false)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume("test", "", false, false, false, false, nil)
	end := make(chan bool)
	routineFinished := make(chan bool)

	for i := 0; i < runtime.NumCPU(); i++ {
		i := i
		go func(messages <-chan amqp.Delivery, end <- chan bool) {
			for {
				select {
				case msg := <-messages:
					log.Printf("msg: %v from worker %v\n", msg, i)
					msg.Ack(false)
					jsonBody, err := json.Marshal(map[string]string{
						"test": "test1",
						"name": "Flo",
					})
					if err != nil {
						log.Println(err)
						continue
					}
					resp, err := http.Post("http://localhost/testing_stuff", "application/json", bytes.NewBuffer(jsonBody))
					if err != nil {
						log.Println(err)
						continue
					}
					defer resp.Body.Close()
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						log.Fatalln(err)
					}
					responseBody := string(body)
					log.Printf("%v, %v", resp.Status, responseBody)
				case <-end:
					log.Println("Finish worker")
					routineFinished <- true
					return
				}
			}
		}(msgs, end)
	}

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	sig := <-gracefulStop
	log.Printf("caught sig %v. Now graceful stop.\n", sig)
	for i := 0; i < runtime.NumCPU(); i++ {
		end <- true
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		<-routineFinished
	}
	log.Println("Finished waiting for the processes. Done!")
}
