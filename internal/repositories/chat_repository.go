package repositories

import (
	"context"

	"backend-AI-Knowledge-Assistant/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepository struct {
	DB *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{DB: db}
}

func (r *ChatRepository) Save(userID, question, answer string) (*models.ChatHistory, error) {
	query := `
		INSERT INTO chat_histories (user_id, question, answer)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, question, answer, created_at
	`

	var chat models.ChatHistory
	err := r.DB.QueryRow(context.Background(), query, userID, question, answer).
		Scan(&chat.ID, &chat.UserID, &chat.Question, &chat.Answer, &chat.CreatedAt)

	return &chat, err
}