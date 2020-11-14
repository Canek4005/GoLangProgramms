package main
import (
    "fmt"
    "net"
    "encoding/json"
)

type User struct {

    Name       string
    Key        string
}
var conn [1000]net.Conn
var key []string
var u []User
var data []byte
var err error
func main() {
  u=make([]User,1000)
  data=make([]byte,1024)
    listener, err0 := net.Listen("tcp", ":8888")
    catchError(err0)

    defer listener.Close()
    fmt.Println("Server is listening...")
    i:=0
    for {

        conn[i], err = listener.Accept()
        if err != nil {
            fmt.Println(err)
            conn[i].Close()
            continue
        }

        n,err2:=conn[i].Read(data)
        catchError(err2)
        er := json.Unmarshal(data[0:n], &u[i])
        catchError(er)
        fmt.Println("Подключен :", u[i].Name)

        i+=1
for k:=0;k<i;k++{
  for j:=0;j<i;j++{
    fmt.Println("пробую ",u[i].Key,"  ",u[j].Key)
        if(u[k].Key==u[j].Key&&k!=j){
        fmt.Println("Сопряжение :",u[k].Name," и ",u[j].Name)
        go handleConnection(k,j)  // запускаем горутину для обработки запроса
        go handleConnection(j,k)
        k=i
        j=i
        break

      }

}

    }
}
}
// обработка подключения
func handleConnection(i int,ide int) {
    defer conn[i].Close()
    greeting(i,ide)
    for {
        // считываем полученные в запросе данные
        input := make([]byte, (1024))
        n, err := conn[i].Read(input)
        if n == 0 || err != nil {
            fmt.Println("Read error:", err)
            break
        }
        source := string(input[0:n])



        // выводим на консоль сервера диагностическую информацию
        fmt.Println(source)
        // отправляем данные клиенту
        conn[ide].Write([]byte(source))
    }
}

func catchError(err error){
  if err != nil {
      fmt.Println(err)
      return
  }
  }
  func greeting(i int,j int){
    salute:="Вы подключены к "+u[j].Name
    conn[i].Write([]byte(salute))


  }
