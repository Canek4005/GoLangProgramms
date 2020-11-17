package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id    int
	Name  string
	Score int
}

var u User

func main() {
	// подключение БД
	DB, err := sql.Open("mysql", "root:2813308004Sesh@/SpaceWander")
	catchError(err)
	defer DB.Close()
	//on Server
	listener, err := net.Listen("tcp", "195.133.196.5:14880")
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

		//запрос на подключение и прием данныхы
		default:
			fmt.Println("Подключен :", u.Name)
			fmt.Println("Его очки  :", u.Score)
			go updateBase(u.Id, u.Name, u.Score, DB)
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
func updateBase(Id int, Name string, Score int, DB *sql.DB) {
	//загрузка данных из БД
	rows, err := DB.Query("select * from SpaceWander.scoreboard")
	catchError(err)
	u := []User{}
	for rows.Next() {
		s := User{}
		err := rows.Scan(&s.Id, &s.Name, &s.Score)
		if err != nil {
			fmt.Println(err)
			continue
		}
		u = append(u, s)
	}
	// пробегаю по массиву и проверяю есть ли это имя
	find := true
	score := true
	for _, p := range u {

		if p.Name == Name {
			find = false

			if p.Score < Score {
				score = false
			}
		}
	}
	// добавление нового пользователя
	if find {
		_, err := DB.Exec("insert into SpaceWander.scoreboard (Id, Name, Score) values (?, ?, ?)",
			Id, Name, Score)
		catchError(err)
	}
	// обновление данных очков
	if !find && !score {
		_, err := DB.Exec("update SpaceWander.scoreboard set Score = ? where Name = ?", Score, Name)
		catchError((err))
	}

	//
}

//отправка таблицы
func sendTable(conn net.Conn, DB *sql.DB) {
	defer conn.Close()
	//получение строки БД
	raws, err := DB.Query("select * from SpaceWander.scoreboard")
	catchError(err)
	users := []User{}
	//обработка строки БД
	for raws.Next() {
		u := User{}
		err := raws.Scan(&u.Id, &u.Name, &u.Score)
		catchError(err)
		users = append(users, u)
	}
	//Отправка БД пользователю
	data_json, err := json.Marshal(users)
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
