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
