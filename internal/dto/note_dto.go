package dto

type CreateNoteRequest struct {
	Title   string `json:"title" validate:"required,min=3"`
	Content string `json:"content"`
}

type UpdateNoteRequest struct {
	Title   *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
}

type NoteResponse struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
