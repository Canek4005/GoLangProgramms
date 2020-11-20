package main

import (
	"fmt"

	"strconv"

	"github.com/go-redis/redis"
)

type User struct {
	Name  string
	Score int
}

func main() {
	fmt.Println("ПОДКЛЮЧАЮС")

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	var u []User
	rows, err := client.HGetAll("SpaceWander").Result()
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("%+v", rows)
	i := 1
	for r, t := range rows {
		var p User
		p.Name = r
		p.Score, _ = strconv.Atoi(t)
		fmt.Println(r, "   ", t)
		i += 1
		u = append(u, p)
	}
	fmt.Println(u)

}
