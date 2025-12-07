package service

import (
    "errors"
    "strings"

    "example.com/notes-api/internal/core"
    "example.com/notes-api/internal/repo"
)

var (
    ErrValidation = errors.New("validation error")
)

type NoteService struct {
    repo repo.NoteRepository
}

func NewNoteService(r repo.NoteRepository) *NoteService {
    return &NoteService{repo: r}
}

func (s *NoteService) CreateNote(title, content string) (*core.Note, error) {
    title = strings.TrimSpace(title)
    if title == "" {
        return nil, ErrValidation
    }

    n := core.Note{
        Title:   title,
        Content: content,
    }
    id, err := s.repo.Create(n)
    if err != nil {
        return nil, err
    }
    return s.repo.GetByID(id)
}

func (s *NoteService) ListNotes() ([]core.Note, error) {
    return s.repo.GetAll()
}

func (s *NoteService) GetNote(id int64) (*core.Note, error) {
    return s.repo.GetByID(id)
}

type NoteUpdateInput struct {
    Title   *string `json:"title"`
    Content *string `json:"content"`
}

func (s *NoteService) UpdateNote(id int64, input NoteUpdateInput) (*core.Note, error) {
    return s.repo.Update(id, func(n *core.Note) error {
        if input.Title != nil {
            title := strings.TrimSpace(*input.Title)
            if title == "" {
                return ErrValidation
            }
            n.Title = title
        }
        if input.Content != nil {
            n.Content = *input.Content
        }
        return nil
    })
}

func (s *NoteService) DeleteNote(id int64) error {
    return s.repo.Delete(id)
}
