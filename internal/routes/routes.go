package routes

import (
	"backend-AI-Knowledge-Assistant/internal/handlers"
	"backend-AI-Knowledge-Assistant/internal/middleware"
	"backend-AI-Knowledge-Assistant/internal/repositories"
	"backend-AI-Knowledge-Assistant/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(r *gin.Engine, db *pgxpool.Pool) {
	userRepo := repositories.NewUserRepository(db)
	noteRepo := repositories.NewNoteRepository(db)
	chatRepo := repositories.NewChatRepository(db)

	authService := services.NewAuthService(userRepo)
	aiService := services.NewAIService()
	noteService := services.NewNoteService(noteRepo, chatRepo, aiService)

	authHandler := handlers.NewAuthHandler(authService)
	noteHandler := handlers.NewNoteHandler(noteService, noteRepo)
	aiHandler := handlers.NewAIHandler(noteService)

	api := r.Group("/api")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
	}

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/notes", noteHandler.GetAll)
		protected.GET("/notes/:id", noteHandler.GetByID)
		protected.POST("/notes", noteHandler.Create)
		protected.PUT("/notes/:id", noteHandler.Update)
		protected.DELETE("/notes/:id", noteHandler.Delete)
		protected.GET("/notes/search", noteHandler.Search)
		protected.GET("/notes/:id/summary", noteHandler.GetSummary)
		protected.GET("/notes/:id/keywords", noteHandler.GetKeywords)

		protected.POST("/ai/ask", aiHandler.Ask)
	}
}