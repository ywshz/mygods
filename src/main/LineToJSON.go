package main

import (
	"github.com/streadway/amqp"
	"os"
	"bufio"
	"io"
	"log"
	"time"
)

func ReadLine(fileName string, handler func([]byte)) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadBytes('\n')
		handler(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://admin:d97aNp@10.120.152.242:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"queue.paidorder", // name
		true, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil, // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ReadLine("/Users/yws/Downloads/result.txt", func(line []byte) {
		err = ch.Publish(
			"", // exchange
			q.Name, // routing key
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        line,
			})
		log.Printf(" [x] Sent %s", line)
		failOnError(err, "Failed to publish a message")
		// 10 msg/second
		time.Sleep(time.Millisecond*100)
	})

}