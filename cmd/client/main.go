package main

import (
	"context"
	"go_redis/internal/client"
	"log"
	"time"
)

func main() {
	c, err := client.New(":3001")
	if err != nil {
		log.Fatalln("Client Error:", err.Error())
	}

	go func() {
		err := c.ReadData()
		if err != nil {
			log.Println("Error reading data:", err.Error())
		}
	}()

	//time.Sleep(time.Second)
	err = c.Put(context.Background(), "name", "hardik")
	if err != nil {
		log.Println("Client Error:", err.Error())
	}
	//time.Sleep(time.Second)
	err = c.Put(context.Background(), "surname", "roongta")
	if err != nil {
		log.Println("Client Error:", err.Error())
	}
	//time.Sleep(time.Second)
	err = c.Put(context.Background(), "city", "guwahati")
	if err != nil {
		log.Println("Client Error:", err.Error())
	}
	//time.Sleep(time.Second)
	err = c.Put(context.Background(), "clg", "iitg")
	if err != nil {
		log.Println("Client Error:", err.Error())
	}
	//time.Sleep(time.Second)
	err = c.Get(context.Background(), "name")
	if err != nil {
		log.Println("Client Error:", err.Error())
	}
	//time.Sleep(time.Second)
	err = c.Put(context.Background(), "color", "red")
	if err != nil {
		log.Println("Client Error:", err.Error())
	}
	//time.Sleep(time.Second)
	err = c.Get(context.Background(), "color")
	if err != nil {
		log.Println("Client Error:", err.Error())
	}
	//time.Sleep(time.Second)
	err = c.Get(context.Background(), "surname")
	if err != nil {
		log.Println("Client Error:", err.Error())
	}
	time.Sleep(time.Second)
	c.Close()
}
