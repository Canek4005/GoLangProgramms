package main

import (
	"fmt"

	"encoding/json"
	"net"
)

type User struct {
	Id    int
	Name  string
	Score int
}

var adress string = ":4541"
var user User

func main() {

	conn, err := net.Dial("tcp", adress)

	if err != nil {
		fmt.Println(err)
		fmt.Println("Отсутствует подключение к серверу")
		return

	} else {
		interactionWithTheServer(conn)
	}

}

//отправка данных о совершенной сессии
func interactionWithTheServer(conn net.Conn) {
	defer conn.Close()
	//Читаю из памяти данные игрока!!!!!!!!!!!!!!!!!!!!!!!!!!
	user.Id = 0
	if user.Name == "" || user.Score == 0 {
		for {
			fmt.Print("Приветствую как тебя зовут?: ")
			_, err := fmt.Scanln(&user.Name)
			if err != nil {
				fmt.Println("Некорректный ввод", err)
				continue
			}
			break
		}
		for {
			fmt.Print("Сколько очков Вы набрали?: ")
			_, err := fmt.Scanln(&user.Score)
			if err != nil {
				fmt.Println("Некорректный ввод", err)
				continue
			}
			break
		}
	}
	json_data, err := json.Marshal(user)
	catchError(err)
	// формирую и отправляю дату
	conn.Write(json_data)
	//
	fmt.Print("Ожидание подключения...")
	// принимаю дату
	buff := make([]byte, 1024)
	n, err := conn.Read(buff)

	fmt.Println(string(buff[0:n]))

}

//запрос на таблицу
func queryTable() {
	//подключение
	conn, err := net.Dial("tcp", adress)
	catchError(err)
	defer conn.Close()
	//запрос на таблицу
	request := User{0, "0", 2813308004}
	data_json, err := json.Marshal(request)
	conn.Write(data_json)
	//получение таблицы и ее обработка
	var data_received []byte
	data_received = make([]byte, 1024*4)
	table := []User{}
	for {
		n, err := conn.Read(data_received)
		catchError(err)
		if err := json.Unmarshal(data_received[0:n], &table); err != nil {
			fmt.Println(err)
			return
		}

	}
	// результат функции
	fmt.Print(table)
}

//обработка ошибок
func catchError(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
