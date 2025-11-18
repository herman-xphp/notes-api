package impl

import (
	"context"
	"notes-api/internal/domain"
	"notes-api/internal/dto"
	"testing"
	"time"

	"gorm.io/gorm"
)

// MockNoteRepository is a mock implementation of NoteRepository
type MockNoteRepository struct {
	notes     map[uint]*domain.Note
	userNotes map[uint][]domain.Note
	createErr error
	updateErr error
	deleteErr error
	findErr   error
	nextID    uint
}

func NewMockNoteRepository() *MockNoteRepository {
	return &MockNoteRepository{
		notes:     make(map[uint]*domain.Note),
		userNotes: make(map[uint][]domain.Note),
		nextID:    1,
	}
}

func (m *MockNoteRepository) Create(ctx context.Context, note *domain.Note) error {
	if m.createErr != nil {
		return m.createErr
	}
	note.ID = m.nextID
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()
	m.nextID++
	m.notes[note.ID] = note
	m.userNotes[note.UserID] = append(m.userNotes[note.UserID], *note)
	return nil
}

func (m *MockNoteRepository) Update(ctx context.Context, note *domain.Note) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, ok := m.notes[note.ID]; !ok {
		return gorm.ErrRecordNotFound
	}
	note.UpdatedAt = time.Now()
	m.notes[note.ID] = note
	return nil
}

func (m *MockNoteRepository) Delete(ctx context.Context, id uint, userID uint) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	note, ok := m.notes[id]
	if !ok || note.UserID != userID {
		return gorm.ErrRecordNotFound
	}
	delete(m.notes, id)
	return nil
}

func (m *MockNoteRepository) FindByID(ctx context.Context, id uint, userID uint) (*domain.Note, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	note, ok := m.notes[id]
	if !ok || note.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	return note, nil
}

func (m *MockNoteRepository) FindAllByUser(ctx context.Context, userID uint, page, limit int) ([]domain.Note, int64, error) {
	if m.findErr != nil {
		return nil, 0, m.findErr
	}
	notes := m.userNotes[userID]
	total := int64(len(notes))

	// Simple pagination
	offset := (page - 1) * limit
	if offset >= len(notes) {
		return []domain.Note{}, total, nil
	}
	end := offset + limit
	if end > len(notes) {
		end = len(notes)
	}

	return notes[offset:end], total, nil
}

func TestNoteService_Create(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockNoteRepository()
	service := NewNoteServiceImpl(mockRepo)

	req := dto.CreateNoteRequest{
		Title:   "Test Note",
		Content: "Test Content",
	}
	userID := uint(1)

	note, err := service.Create(ctx, userID, req)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if note == nil {
		t.Fatal("Create() returned nil note")
	}
	if note.Title != req.Title {
		t.Errorf("Create() Title = %v, want %v", note.Title, req.Title)
	}
	if note.UserID != userID {
		t.Errorf("Create() UserID = %v, want %v", note.UserID, userID)
	}
}

func TestNoteService_GetByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockNoteRepository()
	service := NewNoteServiceImpl(mockRepo)

	// Create a note first
	createReq := dto.CreateNoteRequest{
		Title:   "Test Note",
		Content: "Test Content",
	}
	userID := uint(1)
	createdNote, _ := service.Create(ctx, userID, createReq)

	// Test GetByID
	note, err := service.GetByID(ctx, userID, createdNote.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if note == nil {
		t.Fatal("GetByID() returned nil note")
	}
	if note.ID != createdNote.ID {
		t.Errorf("GetByID() ID = %v, want %v", note.ID, createdNote.ID)
	}

	// Test GetByID with non-existent note
	_, err = service.GetByID(ctx, userID, 999)
	if err == nil {
		t.Error("GetByID() should return error for non-existent note")
	}
}

func TestNoteService_Update(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockNoteRepository()
	service := NewNoteServiceImpl(mockRepo)

	// Create a note first
	createReq := dto.CreateNoteRequest{
		Title:   "Test Note",
		Content: "Test Content",
	}
	userID := uint(1)
	createdNote, _ := service.Create(ctx, userID, createReq)

	// Test Update
	newTitle := "Updated Title"
	updateReq := dto.UpdateNoteRequest{
		Title: &newTitle,
	}
	updatedNote, err := service.Update(ctx, userID, createdNote.ID, updateReq)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updatedNote.Title != newTitle {
		t.Errorf("Update() Title = %v, want %v", updatedNote.Title, newTitle)
	}

	// Test Update with non-existent note
	_, err = service.Update(ctx, userID, 999, updateReq)
	if err == nil {
		t.Error("Update() should return error for non-existent note")
	}
}

func TestNoteService_Delete(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockNoteRepository()
	service := NewNoteServiceImpl(mockRepo)

	// Create a note first
	createReq := dto.CreateNoteRequest{
		Title:   "Test Note",
		Content: "Test Content",
	}
	userID := uint(1)
	createdNote, _ := service.Create(ctx, userID, createReq)

	// Test Delete
	err := service.Delete(ctx, userID, createdNote.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify note is deleted
	_, err = service.GetByID(ctx, userID, createdNote.ID)
	if err == nil {
		t.Error("Delete() should have deleted the note")
	}

	// Test Delete with non-existent note
	err = service.Delete(ctx, userID, 999)
	if err == nil {
		t.Error("Delete() should return error for non-existent note")
	}
}

func TestNoteService_GetAll(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockNoteRepository()
	service := NewNoteServiceImpl(mockRepo)

	userID := uint(1)

	// Create multiple notes
	for i := 0; i < 3; i++ {
		createReq := dto.CreateNoteRequest{
			Title:   "Test Note " + string(rune(i+1)),
			Content: "Test Content",
		}
		service.Create(ctx, userID, createReq)
	}

	// Test GetAll with pagination
	notes, total, err := service.GetAll(ctx, userID, 1, 10)
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}
	if len(notes) != 3 {
		t.Errorf("GetAll() returned %d notes, want 3", len(notes))
	}
	if total != 3 {
		t.Errorf("GetAll() total = %d, want 3", total)
	}
}

func TestNoteService_GetAll_Pagination(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockNoteRepository()
	service := NewNoteServiceImpl(mockRepo)

	userID := uint(1)

	// Create 25 notes for pagination testing
	totalNotes := 25
	for i := 0; i < totalNotes; i++ {
		createReq := dto.CreateNoteRequest{
			Title:   "Test Note " + string(rune(i+1)),
			Content: "Test Content " + string(rune(i+1)),
		}
		service.Create(ctx, userID, createReq)
	}

	// Test pagination: page 1, limit 10
	notes, total, err := service.GetAll(ctx, userID, 1, 10)
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}
	if len(notes) != 10 {
		t.Errorf("GetAll() page 1 returned %d notes, want 10", len(notes))
	}
	if total != int64(totalNotes) {
		t.Errorf("GetAll() total = %d, want %d", total, totalNotes)
	}

	// Test pagination: page 2, limit 10
	notes2, total2, err := service.GetAll(ctx, userID, 2, 10)
	if err != nil {
		t.Fatalf("GetAll() page 2 error = %v", err)
	}
	if len(notes2) != 10 {
		t.Errorf("GetAll() page 2 returned %d notes, want 10", len(notes2))
	}
	if total2 != int64(totalNotes) {
		t.Errorf("GetAll() page 2 total = %d, want %d", total2, totalNotes)
	}

	// Test pagination: page 3, limit 10 (should have 5 notes)
	notes3, total3, err := service.GetAll(ctx, userID, 3, 10)
	if err != nil {
		t.Fatalf("GetAll() page 3 error = %v", err)
	}
	if len(notes3) != 5 {
		t.Errorf("GetAll() page 3 returned %d notes, want 5", len(notes3))
	}
	if total3 != int64(totalNotes) {
		t.Errorf("GetAll() page 3 total = %d, want %d", total3, totalNotes)
	}

	// Test pagination: page 4, limit 10 (should be empty)
	notes4, total4, err := service.GetAll(ctx, userID, 4, 10)
	if err != nil {
		t.Fatalf("GetAll() page 4 error = %v", err)
	}
	if len(notes4) != 0 {
		t.Errorf("GetAll() page 4 returned %d notes, want 0", len(notes4))
	}
	if total4 != int64(totalNotes) {
		t.Errorf("GetAll() page 4 total = %d, want %d", total4, totalNotes)
	}

	// Test pagination: page 1, limit 20
	notes5, total5, err := service.GetAll(ctx, userID, 1, 20)
	if err != nil {
		t.Fatalf("GetAll() limit 20 error = %v", err)
	}
	if len(notes5) != 20 {
		t.Errorf("GetAll() limit 20 returned %d notes, want 20", len(notes5))
	}
	if total5 != int64(totalNotes) {
		t.Errorf("GetAll() limit 20 total = %d, want %d", total5, totalNotes)
	}

	// Test pagination: page 2, limit 20 (should have 5 notes)
	notes6, total6, err := service.GetAll(ctx, userID, 2, 20)
	if err != nil {
		t.Fatalf("GetAll() page 2 limit 20 error = %v", err)
	}
	if len(notes6) != 5 {
		t.Errorf("GetAll() page 2 limit 20 returned %d notes, want 5", len(notes6))
	}
	if total6 != int64(totalNotes) {
		t.Errorf("GetAll() page 2 limit 20 total = %d, want %d", total6, totalNotes)
	}

	// Test pagination: page 1, limit 25 (all notes)
	notes7, total7, err := service.GetAll(ctx, userID, 1, 25)
	if err != nil {
		t.Fatalf("GetAll() limit 25 error = %v", err)
	}
	if len(notes7) != 25 {
		t.Errorf("GetAll() limit 25 returned %d notes, want 25", len(notes7))
	}
	if total7 != int64(totalNotes) {
		t.Errorf("GetAll() limit 25 total = %d, want %d", total7, totalNotes)
	}

	// Test pagination: page 1, limit 100 (more than total, should return all)
	notes8, total8, err := service.GetAll(ctx, userID, 1, 100)
	if err != nil {
		t.Fatalf("GetAll() limit 100 error = %v", err)
	}
	if len(notes8) != totalNotes {
		t.Errorf("GetAll() limit 100 returned %d notes, want %d", len(notes8), totalNotes)
	}
	if total8 != int64(totalNotes) {
		t.Errorf("GetAll() limit 100 total = %d, want %d", total8, totalNotes)
	}

	// Verify notes are different between pages
	if notes[0].ID == notes2[0].ID {
		t.Error("GetAll() page 1 and page 2 should return different notes")
	}
}

func TestNoteService_GetAll_Pagination_EdgeCases(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockNoteRepository()
	service := NewNoteServiceImpl(mockRepo)

	userID := uint(1)

	// Create 5 notes
	for i := 0; i < 5; i++ {
		createReq := dto.CreateNoteRequest{
			Title:   "Test Note " + string(rune(i+1)),
			Content: "Test Content",
		}
		service.Create(ctx, userID, createReq)
	}

	// Test page 0 (should default to page 1)
	notes, total, err := service.GetAll(ctx, userID, 0, 10)
	if err != nil {
		t.Fatalf("GetAll() page 0 error = %v", err)
	}
	if len(notes) != 5 {
		t.Errorf("GetAll() page 0 returned %d notes, want 5", len(notes))
	}
	if total != 5 {
		t.Errorf("GetAll() page 0 total = %d, want 5", total)
	}

	// Test negative page (should default to page 1)
	notes2, _, err := service.GetAll(ctx, userID, -1, 10)
	if err != nil {
		t.Fatalf("GetAll() negative page error = %v", err)
	}
	if len(notes2) != 5 {
		t.Errorf("GetAll() negative page returned %d notes, want 5", len(notes2))
	}

	// Test limit 0 (should default to limit 10)
	notes3, total3, err := service.GetAll(ctx, userID, 1, 0)
	if err != nil {
		t.Fatalf("GetAll() limit 0 error = %v", err)
	}
	if len(notes3) != 5 {
		t.Errorf("GetAll() limit 0 returned %d notes, want 5", len(notes3))
	}
	if total3 != 5 {
		t.Errorf("GetAll() limit 0 total = %d, want 5", total3)
	}

	// Test empty result (no notes)
	userID2 := uint(2)
	notes4, total4, err := service.GetAll(ctx, userID2, 1, 10)
	if err != nil {
		t.Fatalf("GetAll() empty result error = %v", err)
	}
	if len(notes4) != 0 {
		t.Errorf("GetAll() empty result returned %d notes, want 0", len(notes4))
	}
	if total4 != 0 {
		t.Errorf("GetAll() empty result total = %d, want 0", total4)
	}
}
