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
