package repositories

import (
	"context"

	"backend-AI-Knowledge-Assistant/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NoteRepository struct {
	DB *pgxpool.Pool
}

func NewNoteRepository(db *pgxpool.Pool) *NoteRepository {
	return &NoteRepository{DB: db}
}

func (r *NoteRepository) Create(userID, title, content string) (*models.Note, error) {
	query := `
		INSERT INTO notes (user_id, title, content, ai_status)
		VALUES ($1, $2, $3, 'pending')
		RETURNING id, user_id, title, content, ai_status, created_at, updated_at
	`

	var note models.Note
	err := r.DB.QueryRow(context.Background(), query, userID, title, content).
		Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.AIStatus, &note.CreatedAt, &note.UpdatedAt)

	return &note, err
}

func (r *NoteRepository) FindAll(userID string) ([]models.Note, error) {
	query := `
		SELECT id, user_id, title, content, ai_status, created_at, updated_at
		FROM notes
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`

	rows, err := r.DB.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.AIStatus, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (r *NoteRepository) FindByID(noteID, userID string) (*models.Note, error) {
	query := `
		SELECT id, user_id, title, content, ai_status, created_at, updated_at
		FROM notes
		WHERE id = $1 AND user_id = $2
	`

	var note models.Note
	err := r.DB.QueryRow(context.Background(), query, noteID, userID).
		Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.AIStatus, &note.CreatedAt, &note.UpdatedAt)

	return &note, err
}

func (r *NoteRepository) Update(noteID, userID, title, content string) (*models.Note, error) {
	query := `
		UPDATE notes
		SET title = $1, content = $2, ai_status = 'pending', updated_at = NOW()
		WHERE id = $3 AND user_id = $4
		RETURNING id, user_id, title, content, ai_status, created_at, updated_at
	`

	var note models.Note
	err := r.DB.QueryRow(context.Background(), query, title, content, noteID, userID).
		Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.AIStatus, &note.CreatedAt, &note.UpdatedAt)

	return &note, err
}

func (r *NoteRepository) Delete(noteID, userID string) error {
	query := `DELETE FROM notes WHERE id = $1 AND user_id = $2`
	_, err := r.DB.Exec(context.Background(), query, noteID, userID)
	return err
}

func (r *NoteRepository) Search(userID, keyword string) ([]models.Note, error) {
	query := `
		SELECT id, user_id, title, content, ai_status, created_at, updated_at
		FROM notes
		WHERE user_id = $1
		AND to_tsvector('simple', coalesce(title,'') || ' ' || coalesce(content,''))
		@@ plainto_tsquery('simple', $2)
		ORDER BY updated_at DESC
	`

	rows, err := r.DB.Query(context.Background(), query, userID, keyword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.AIStatus, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (r *NoteRepository) SaveSummary(noteID, summary, model string) error {
	query := `
		INSERT INTO note_summaries (note_id, summary_text, model_used)
		VALUES ($1, $2, $3)
		ON CONFLICT (note_id)
		DO UPDATE SET
			summary_text = EXCLUDED.summary_text,
			model_used = EXCLUDED.model_used,
			updated_at = NOW()
	`
	_, err := r.DB.Exec(context.Background(), query, noteID, summary, model)
	return err
}

func (r *NoteRepository) GetSummary(noteID string) (*models.NoteSummary, error) {
	query := `
		SELECT id, note_id, summary_text, model_used, created_at, updated_at
		FROM note_summaries
		WHERE note_id = $1
	`

	var summary models.NoteSummary
	err := r.DB.QueryRow(context.Background(), query, noteID).
		Scan(&summary.ID, &summary.NoteID, &summary.SummaryText, &summary.ModelUsed, &summary.CreatedAt, &summary.UpdatedAt)

	return &summary, err
}

func (r *NoteRepository) ReplaceKeywords(noteID string, keywords []models.NoteKeyword) error {
	tx, err := r.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `DELETE FROM note_keywords WHERE note_id = $1`, noteID)
	if err != nil {
		return err
	}

	for _, k := range keywords {
		_, err = tx.Exec(
			context.Background(),
			`INSERT INTO note_keywords (note_id, keyword, score) VALUES ($1, $2, $3)`,
			noteID, k.Keyword, k.Score,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}

func (r *NoteRepository) GetKeywords(noteID string) ([]models.NoteKeyword, error) {
	query := `
		SELECT id, note_id, keyword, score, created_at
		FROM note_keywords
		WHERE note_id = $1
		ORDER BY score DESC, keyword ASC
	`

	rows, err := r.DB.Query(context.Background(), query, noteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keywords []models.NoteKeyword
	for rows.Next() {
		var item models.NoteKeyword
		if err := rows.Scan(&item.ID, &item.NoteID, &item.Keyword, &item.Score, &item.CreatedAt); err != nil {
			return nil, err
		}
		keywords = append(keywords, item)
	}

	return keywords, nil
}