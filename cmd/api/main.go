package main

import (
    "log"
    "net/http"

    httpx "example.com/notes-api/internal/http"
    "example.com/notes-api/internal/http/handlers"
    "example.com/notes-api/internal/core/service"
    "example.com/notes-api/internal/repo"
)

func main() {
    // Инициализация репозитория и сервиса
    rp := repo.NewNoteRepoMem()
    svc := service.NewNoteService(rp)
    h := handlers.NewHandler(svc)

    router := httpx.NewRouter(h)

    addr := ":8080" // слушаем на 0.0.0.0:8080
    log.Println("Server started at", addr)
    log.Fatal(http.ListenAndServe(addr, router))
}
