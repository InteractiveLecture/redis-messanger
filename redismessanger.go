package redismessanger

import (
	"fmt"
	"os"

	"gopkg.in/redis.v3"
)

type Messanger struct {
	client *redis.Client
}

func (m *Messanger) publish(channel, message string) error {
	pubsub, err := m.client.Subscribe(channel)
	if err != nil {
		panic(err)
	}
	defer pubsub.Close()
	multi := client.Multi()
	multi.LPush(channel, message)
	multi.LTrim(channel, 0, 50)
	_, err = multi.Exec(func() error { return nil })
	if err != nil {
		return err
	}
	return client.Publish(channel, message).Err()
}

type MessageListener func(message string)

func (m *Messanger) Subscribe(channel string, callback MessageListener) error {
	pubsub, err := m.client.Subscribe(channel)
	if err != nil {
		return err
	}

}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDISHOST") + ":6379",
		Password: "",
		DB:       0,
	})
	pubsub, err := client.Subscribe("authentication-service:new-users")
	if err != nil {
		panic(err)
	}
	defer pubsub.Close()
	multi := client.Multi()
	multi.LPush("authentication-service:new-users", "123")
	multi.LTrim("authentication-service:new-users", 0, 5)
	_, err = multi.Exec(func() error { return nil })
	if err != nil {
		panic(err)
	}
	err = client.Publish("authentication-service:new-users", "t").Err()
	if err != nil {
		panic(err)
	}
	msg, err := pubsub.ReceiveMessage()
	if err != nil {
		panic(err)
	}
	fmt.Println(msg)

	result, err := client.LRange("authentication-service:new-users", 0, 5).Result()
	if err != nil {
		panic(err)
	}
	for _, v := range result {
		fmt.Println(v)
	}
}
