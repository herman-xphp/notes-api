package dto

type CreateNoteRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content"`
}

type UpdateNoteRequest struct {
	Title   *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
}
