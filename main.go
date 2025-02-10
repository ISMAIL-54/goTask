package main

import (
    "log"
)

func main() {
    storage, err := NewPostgresDB()
    if err != nil {
        log.Fatal(err)
    }
    err = storage.Init()
    if err != nil {
        log.Fatal(err)
    }
    server := NewAPIServer(":8080", storage)
    server.Run()
}
