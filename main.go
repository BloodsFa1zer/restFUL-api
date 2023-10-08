package main

import (
	"app3.1/database"
	"fmt"
	"github.com/go-redis/redis"
)

func main() {
	//e := echo.New()
	//e.HidePort = true
	//
	//Middleware.UserAuth(e)
	//routes.UserRoute(e)
	//
	//e.Logger.Fatal(e.Start(":6000"))

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	// we can call set with a Key' and a 'Value'
	err = client.Set("1", database.User{}, 0).Err()
	// if there has been an error setting the value
	// handle the error
	if err != nil {
		fmt.Println(err)
	}

	val, err := client.Get("name").Result()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(val)
}
