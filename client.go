package main
import(
  "fmt"

  "net"
  "encoding/json"
  )
  type User struct {

      Name       string
      Key        string
  }
var user User
func main(){
for{
  fmt.Print("Приветствую как тебя зовут?: ")
  _, err := fmt.Scanln(&user.Name)
  if err != nil {
      fmt.Println("Некорректный ввод", err)
      continue
  }
  break
}
for{
  fmt.Print("Какой ключ Вы будете использовать?: ")
  _, err := fmt.Scanln(&user.Key)
  if err != nil {
      fmt.Println("Некорректный ввод", err)
      continue
  }
  break
}

conn, err := net.Dial("tcp", ":8888")
catchError(err)
defer conn.Close()

json_data, err1 := json.Marshal(user)
catchError(err1)

conn.Write(json_data)
fmt.Print("Ожидание подключения...")
buff := make([]byte, 1024)
n, err := conn.Read(buff)

fmt.Println("  ",string(buff[0:n]))
  for{
      go sendMessage(conn)
      buff := make([]byte, 1024)
      n, err := conn.Read(buff)
      if err !=nil{ break}
      fmt.Println("Получено сообщение :",string(buff[0:n]))

  }
  }
  func sendMessage(conn net.Conn){
    for{
      var source string
      _, err := fmt.Scanln(&source)
      if err != nil {
          fmt.Println("Некорректный ввод", err)
          continue
      }
      conn.Write([]byte(source))

  }
}
  func catchError(err error){
    if err != nil {
        fmt.Println(err)
        return
    }
  }
