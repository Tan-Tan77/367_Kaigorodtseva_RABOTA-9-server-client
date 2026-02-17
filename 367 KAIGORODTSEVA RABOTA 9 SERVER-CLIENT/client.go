package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("Не удалось подключиться к серверу:", err)
        return
    }
    defer conn.Close()

    fmt.Println("Подключено к серверу. Для выхода введите 'exit'.")

 
    messageChan := make(chan string, 5)
    

    go displayMessages(messageChan)

  
    go func() {
        reader := bufio.NewReader(conn)
        for {
            message, err := reader.ReadString('\n')
            if err != nil {
                fmt.Println("Соединение с сервером разорвано")
                os.Exit(0)
            }
            messageChan <- message
        }
    }()

 
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        text := scanner.Text()
        
        if text == "exit" {
            fmt.Println("Выход...")
            return
        }


        _, err := fmt.Fprintln(conn, text)
        if err != nil {
            fmt.Println("Ошибка отправки сообщения:", err)
            return
        }
    }

   
    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
    <-signalChan
    fmt.Println("\nЗавершение работы...")
}

func displayMessages(messageChan chan string) {
    for message := range messageChan {
        fmt.Print(message)
    }
}