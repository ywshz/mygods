package main

import (
	"github.com/streadway/amqp"
	"os"
	"bufio"
	"io"
	"log"
	"time"
	"io/ioutil"
	"github.com/henrylee2cn/pholcus/common/simplejson"
	"fmt"
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
	configData,_ := ioutil.ReadFile("./config.json")
	fmt.Println(string(configData))
	configJson,_ := simplejson.NewJson(configData)
	conn ,err := amqp.Dial(configJson.Get("conn").MustString())
	file := configJson.Get("file").MustString()
	exchange := configJson.Get("exchange").MustString()
	routekey := configJson.Get("routekey").MustString()
	interval := configJson.Get("sendInterval").MustInt64()

	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	//q, err := ch.QueueDeclare(
	//	queue, // name
	//	true, // durable
	//	false, // delete when unused
	//	false, // exclusive
	//	false, // no-wait
	//	nil, // arguments
	//)
	//q, err := ch.QueueDeclare(
	//	"queue.paidorder", // name
	//	true, // durable
	//	false, // delete when unused
	//	false, // exclusive
	//	false, // no-wait
	//	nil, // arguments
	//)
	failOnError(err, "Failed to declare a queue")

	ReadLine(file, func(line []byte) {
		if string(line) == "" {
			return
		}
		err = ch.Publish(
			exchange, // exchange
			routekey, // routing key
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        line,
			})
		log.Printf(" [x] Sent %s", line)
		failOnError(err, "Failed to publish a message")
		// 1000/num 条/秒
		time.Sleep(time.Millisecond * time.Duration(interval))
	})

}