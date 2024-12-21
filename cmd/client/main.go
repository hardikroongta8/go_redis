package main

import (
	"context"
	"fmt"
	"go_redis/internal/client"
	"log"
)

func main() {
	c, err := client.New(":3001")
	if err != nil {
		log.Fatalln("Client Error:", err.Error())
	}
	c.WG.Add(1)
	go func() {
		err := c.ReadData()
		if err != nil {
			log.Println("Error reading data:", err.Error())
		}
		c.WG.Done()
	}()

	err = c.Put(context.Background(), "name", fmt.Sprintf("Hardik"))
	if err != nil {
		log.Println("Client Error:", err.Error())
	}

	err = c.Put(context.Background(), "surname", fmt.Sprintf("Roongta"))
	if err != nil {
		log.Println("Client Error:", err.Error())
	}

	err = c.Put(context.Background(), "city", fmt.Sprintf("Guwahati"))
	if err != nil {
		log.Println("Client Error:", err.Error())
	}

	err = c.Put(context.Background(), "clg", "iitg")
	if err != nil {
		log.Println("Client Error:", err.Error())
	}

	err = c.Get(context.Background(), "name")
	if err != nil {
		log.Println("Client Error:", err.Error())
	}

	err = c.Put(context.Background(), "color", fmt.Sprintf("Red"))
	if err != nil {
		log.Println("Client Error:", err.Error())
	}

	err = c.Get(context.Background(), "color")
	if err != nil {
		log.Println("Client Error:", err.Error())
	}

	err = c.Get(context.Background(), "surname")
	if err != nil {
		log.Println("Client Error:", err.Error())
	}

	c.WG.Wait()
}
