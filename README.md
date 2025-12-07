
           /\
       )  ( ')
      (  /  )
       \(__)|
# практическое задание 11
## Шишков А.Д. ЭФМО-02-21
## Тема
Проектирование REST API (CRUD для заметок). Разработка структуры 
## Цели
- Освоить принципы проектирования REST API.
- Научиться разрабатывать структуру проекта backend-приложения на Go.
- Спроектировать и реализовать CRUD-интерфейс (Create, Read, Update, Delete) для сущности «Заметка».
- Освоить применение слоистой архитектуры (handler → service → repository).
- Подготовить основу для интеграции с базой данных и JWT-аутентификацией в следующих занятиях.
## Описание проекта
В рамках задания требуется разработать серверное приложение на языке Go, реализующее REST API для управления сущностью «заметка» (Note).
Необходимо создать:
- доменную модель заметки;
- слой репозитория с хранением данных в оперативной памяти (in-memory);
- слой бизнес-логики (service) с валидацией данных;
- HTTP-обработчики для обработки запросов;
- маршруты API с использованием фреймворка chi;
  
Полноценный CRUD:
- создание заметки (POST),
- получение всех заметок (GET),
- получение заметки по ID (GET),
- обновление заметки (PATCH),
- удаление заметки (DELETE).
После создания API его требуется протестировать с помощью Postman, сформировав запросы для всех операций согласно структуре маршрутов.
### Теоретические положения REST API и CRUD

### REST API

**REST (Representational State Transfer)** — архитектурный стиль взаимодействия компонентов распределённого приложения в сети. REST-архитектура основывается на следующих принципах:

1. **Клиент-серверная архитектура** — разделение ответственности между клиентом (пользовательский интерфейс) и сервером (хранение данных).

2. **Отсутствие состояния (Stateless)** — каждый запрос от клиента содержит всю необходимую информацию для его обработки. Сервер не хранит состояние клиента между запросами.

3. **Кэширование** — ответы сервера могут быть помечены как кэшируемые или некэшируемые.

4. **Единообразие интерфейса** — использование стандартных HTTP-методов и унифицированных URL для доступа к ресурсам.

5. **Многоуровневая система** — клиент не может определить, подключён ли он напрямую к серверу или через промежуточные узлы.

### CRUD операции

**CRUD** — акроним, обозначающий четыре базовые операции над данными:

| Операция   | HTTP-метод |        Описание         | Пример URL                                       |
|------------|------------|-------------------------|--------------------------------------------------|
| **C**reate | POST       | Создание нового ресурса | `POST /api/v1/notes`                             |
| **R**ead   | GET        | Чтение ресурса(ов)      | `GET /api/v1/notes` или `GET /api/v1/notes/{id}` |
| **U**pdate | PUT/PATCH  | Обновление ресурса      | `PATCH /api/v1/notes/{id}`                       |
| **D**elete | DELETE     | Удаление ресурса        | `DELETE /api/v1/notes/{id}`                      |

### HTTP-коды ответов, используемые в проекте

| Код |  Название   | Описание                            |
|-----|-------------|-------------------------------------|
| 200 | OK          | Успешный запрос                     |
| 201 | Created     | Ресурс успешно создан               |
| 204 | No Content  | Успешное удаление (без тела ответа) |
| 400 | Bad Request | Некорректный запрос                 |
| 404 | Not Found   | Ресурс не найден                    |
| 500 | Internal Server Error | Внутренняя ошибка сервера|

---
### 3. Структура созданного проекта 

<img width="303" height="472" alt="image" src="https://github.com/user-attachments/assets/257c4bda-2c8f-47cd-b5cb-64080b4f6a4b" /> 

### Примеры кода основных файлов

### main.go — Точка входа в приложение

```go
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

    addr := ":8080"
    log.Println("Server started at", addr)
    log.Fatal(http.ListenAndServe(addr, router))
}
```

**Описание:** Файл `main.go` является точкой входа в приложение. Здесь происходит инициализация всех компонентов системы по принципу Dependency Injection: создаётся репозиторий, сервис и обработчики. Затем настраивается маршрутизатор и запускается HTTP-сервер на порту 8080.

---

### note_mem.go — In-memory репозиторий

```go
package repo

import (
    "errors"
    "sync"
    "time"

    "example.com/notes-api/internal/core"
)

var (
    ErrNoteNotFound = errors.New("note not found")
)

// NoteRepository — интерфейс репозитория
type NoteRepository interface {
    Create(note core.Note) (int64, error)
    GetAll() ([]core.Note, error)
    GetByID(id int64) (*core.Note, error)
    Update(id int64, updateFn func(*core.Note) error) (*core.Note, error)
    Delete(id int64) error
}

// NoteRepoMem — in-memory реализация
type NoteRepoMem struct {
    mu    sync.RWMutex
    notes map[int64]*core.Note
    next  int64
}

func NewNoteRepoMem() *NoteRepoMem {
    return &NoteRepoMem{
        notes: make(map[int64]*core.Note),
    }
}

func (r *NoteRepoMem) Create(n core.Note) (int64, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.next++
    n.ID = r.next
    now := time.Now().UTC()
    n.CreatedAt = now
    n.UpdatedAt = nil

    r.notes[n.ID] = &n
    return n.ID, nil
}

func (r *NoteRepoMem) GetAll() ([]core.Note, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    result := make([]core.Note, 0, len(r.notes))
    for _, n := range r.notes {
        result = append(result, *n)
    }
    return result, nil
}

func (r *NoteRepoMem) GetByID(id int64) (*core.Note, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    n, ok := r.notes[id]
    if !ok {
        return nil, ErrNoteNotFound
    }
    copy := *n
    return &copy, nil
}

func (r *NoteRepoMem) Update(id int64, updateFn func(*core.Note) error) (*core.Note, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    n, ok := r.notes[id]
    if !ok {
        return nil, ErrNoteNotFound
    }

    if err := updateFn(n); err != nil {
        return nil, err
    }
    now := time.Now().UTC()
    n.UpdatedAt = &now

    copy := *n
    return &copy, nil
}

func (r *NoteRepoMem) Delete(id int64) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, ok := r.notes[id]; !ok {
        return ErrNoteNotFound
    }
    delete(r.notes, id)
    return nil
}
```

**Описание:** Файл `note_mem.go` содержит реализацию репозитория для хранения заметок в оперативной памяти. Используется `sync.RWMutex` для обеспечения потокобезопасности при конкурентном доступе к данным. Интерфейс `NoteRepository` позволяет легко заменить реализацию на базу данных без изменения бизнес-логики.

---

### handlers/notes.go — HTTP-обработчики

```go
package handlers

import (
    "encoding/json"
    "errors"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"

    "example.com/notes-api/internal/core/service"
    "example.com/notes-api/internal/repo"
)

type Handler struct {
    Service *service.NoteService
}

func NewHandler(s *service.NoteService) *Handler {
    return &Handler{Service: s}
}

// вспомогательная функция для ошибок
func writeError(w http.ResponseWriter, status int, msg string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(map[string]string{
        "error": msg,
    })
}

// POST /api/v1/notes
func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Title   string `json:"title"`
        Content string `json:"content"`
    }

    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        writeError(w, http.StatusBadRequest, "invalid JSON")
        return
    }

    note, err := h.Service.CreateNote(input.Title, input.Content)
    if err != nil {
        if errors.Is(err, service.ErrValidation) {
            writeError(w, http.StatusBadRequest, "title is required")
            return
        }
        writeError(w, http.StatusInternalServerError, "internal error")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated) // 201
    _ = json.NewEncoder(w).Encode(note)
}

// GET /api/v1/notes
func (h *Handler) ListNotes(w http.ResponseWriter, r *http.Request) {
    notes, err := h.Service.ListNotes()
    if err != nil {
        writeError(w, http.StatusInternalServerError, "internal error")
        return
    }
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(notes)
}

// GET /api/v1/notes/{id}
func (h *Handler) GetNote(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid id")
        return
    }

    note, err := h.Service.GetNote(id)
    if err != nil {
        if errors.Is(err, repo.ErrNoteNotFound) {
            writeError(w, http.StatusNotFound, "note not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "internal error")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(note)
}

// PATCH /api/v1/notes/{id}
func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid id")
        return
    }

    var input service.NoteUpdateInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        writeError(w, http.StatusBadRequest, "invalid JSON")
        return
    }

    note, err := h.Service.UpdateNote(id, input)
    if err != nil {
        if errors.Is(err, repo.ErrNoteNotFound) {
            writeError(w, http.StatusNotFound, "note not found")
            return
        }
        if errors.Is(err, service.ErrValidation) {
            writeError(w, http.StatusBadRequest, "invalid data")
            return
        }
        writeError(w, http.StatusInternalServerError, "internal error")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(note)
}

// DELETE /api/v1/notes/{id}
func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid id")
        return
    }

    if err := h.Service.DeleteNote(id); err != nil {
        if errors.Is(err, repo.ErrNoteNotFound) {
            writeError(w, http.StatusNotFound, "note not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "internal error")
        return
    }

    w.WriteHeader(http.StatusNoContent) // 204, без тела
}
```

**Описание:** Файл `notes.go` содержит HTTP-обработчики для всех CRUD-операций. Каждый обработчик:
1. Извлекает данные из запроса (параметры URL, тело запроса)
2. Вызывает соответствующий метод сервиса
3. Формирует HTTP-ответ в формате JSON

---
