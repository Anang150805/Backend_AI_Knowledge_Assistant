package models

import "time"

type Note struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AIStatus  string    `json:"ai_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NoteSummary struct {
	ID          string    `json:"id"`
	NoteID       string    `json:"note_id"`
	SummaryText  string    `json:"summary_text"`
	ModelUsed    string    `json:"model_used"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type NoteKeyword struct {
	ID        string    `json:"id"`
	NoteID    string    `json:"note_id"`
	Keyword   string    `json:"keyword"`
	Score     float64   `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}