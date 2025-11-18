package dto

import "time"

// NoteResponse adalah response DTO untuk Note (tanpa field internal)
type NoteResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToNoteResponse mengkonversi domain.Note ke NoteResponse
func ToNoteResponse(id uint, title, content string, createdAt, updatedAt time.Time) NoteResponse {
	return NoteResponse{
		ID:        id,
		Title:     title,
		Content:   content,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

