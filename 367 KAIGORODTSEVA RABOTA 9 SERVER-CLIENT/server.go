package main

import (
    "bufio"
    "fmt"
    "net"
    "strings"
    "sync"
)

type Client struct {
    conn net.Conn
    name string
}

var (
    clients     = make(map[*Client]bool)
    clientsLock sync.RWMutex
    messageChan = make(chan string, 10)
)

func main() {
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println("Ошибка запуска сервера:", err)
        return
    }
    defer listener.Close()

    fmt.Println("Сервер запущен на порту 8080")

   
    go displayMessages()

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Ошибка подключения:", err)
            continue
        }

        
        go handleClient(conn)
    }
}

func handleClient(conn net.Conn) {
    defer conn.Close()

    conn.Write([]byte("Введите ваше имя: "))
    reader := bufio.NewReader(conn)
    name, _ := reader.ReadString('\n')
    name = strings.TrimSpace(name)

    client := &Client{conn: conn, name: name}

  
    clientsLock.Lock()
    clients[client] = true
    clientsLock.Unlock()

    messageChan <- fmt.Sprintf("Пользователь %s присоединился к чату", name)

 
    conn.Write([]byte(fmt.Sprintf("Добро пожаловать в чат, %s!\n", name)))


    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            break
        }

        message = strings.TrimSpace(message)
        if message == "" {
            continue
        }

      
        broadcastMessage(fmt.Sprintf("%s: %s", name, message))
    }


    clientsLock.Lock()
    delete(clients, client)
    clientsLock.Unlock()

    messageChan <- fmt.Sprintf("Пользователь %s покинул чат", name)
}

func broadcastMessage(message string) {
    clientsLock.RLock()
    defer clientsLock.RUnlock()

  
    for client := range clients {
        _, err := client.conn.Write([]byte(message + "\n"))
        if err != nil {
          
        }
    }

  
    messageChan <- message
}

func displayMessages() {
    for message := range messageChan {
        fmt.Println(message)
    }
}