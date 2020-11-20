package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/go-redis/redis"
)

type User struct {
	Name  string
	Score int
}

var u User

func main() {
	// подключение БД
	DB := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	//on Server
	listener, err := net.Listen("tcp", ":1480")
	catchError(err)
	defer listener.Close()
	//
	fmt.Println("Server is listening...")

	for {
		var data []byte
		data = make([]byte, 1024)
		// accept conn
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		//catch data
		n, err := conn.Read(data)
		catchError(err)
		//translate json
		fmt.Println("Расшифровываю json")
		if err := json.Unmarshal(data[0:n], &u); err != nil {
			fmt.Println(err)
			return
		}

		switch u.Score {
		//запрос на таблицу
		case 281330800:
			sendTable(conn, DB)

		//запрос на подключение и прием данных
		default:
			fmt.Println("Подключен :", u.Name)
			fmt.Println("Его очки  :", u.Score)
			go updateBase(u.Name, u.Score, DB)
			go handleConnection(conn)

		}
	}
}

// обработка подключения
func handleConnection(conn net.Conn) {
	defer conn.Close()
	salute := "Вы подключены"
	conn.Write([]byte(salute))
	fmt.Println("Отправлено приветствие")

	conn.Close()
}
func updateBase(Name string, Score int, DB *redis.Client) {

	// ищу имя в БД

	isFind, _ := DB.HExists("SpaceWander", Name).Result()

	// добавление нового пользователя
	if !isFind {
		DB.HSet("SpaceWander", Name, Score)
	}
	// обновление данных очков
	if isFind {
		val, _ := DB.HGet("SpaceWander", Name).Result()
		value, _ := strconv.Atoi(val)
		if value < Score {

			DB.HSet("SpaceWander", Name, Score)
		}
	}

}

//отправка таблицы
func sendTable(conn net.Conn, DB *redis.Client) {
	defer conn.Close()

	var users []User
	rows, err := DB.HGetAll("SpaceWander").Result()
	if err != nil {
		fmt.Print(err)
	}

	i := 1
	for r, t := range rows {
		var p User
		p.Name = r
		p.Score, _ = strconv.Atoi(t)
		fmt.Println(r, "   ", t)
		i += 1
		users = append(users, p)
	}
	// Сортировка пользователей

	var k User
	for i := 0; i < len(users); i++ {
		for j := 0; j < len(users)-1; j++ {
			if users[j].Score < users[j+1].Score {
				k = users[j]
				users[j] = users[j+1]
				users[j+1] = k
			}
		}
	}
	fmt.Println(users)
	// отправляю первых 9 человек
	var cutUsers []User
	if len(users) > 9 {
		cutUsers = users[:9]

	} else {
		cutUsers = users
	}
	//Отправка БД пользователю
	data_json, err := json.Marshal(cutUsers)
	conn.Write(data_json)
	fmt.Println("Отправил таблицу")
	fmt.Println(users)
	conn.Close()
}

//обработка ошибок
func catchError(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
