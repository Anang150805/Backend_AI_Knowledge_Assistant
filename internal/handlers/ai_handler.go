package handlers

import (
	"net/http"
	"strings"

	"backend-AI-Knowledge-Assistant/internal/services"
	"backend-AI-Knowledge-Assistant/pkg/response"

	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	NoteService *services.NoteService
}

func NewAIHandler(noteService *services.NoteService) *AIHandler {
	return &AIHandler{NoteService: noteService}
}

type askRequest struct {
	Question string `json:"question"`
}

func (h *AIHandler) Ask(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	var req askRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	if strings.TrimSpace(req.Question) == "" {
		response.Error(c, http.StatusBadRequest, "question wajib diisi")
		return
	}

	answer, err := h.NoteService.AskAI(userID, req.Question)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "ai answer fetched", gin.H{
		"answer": answer,
	})
}