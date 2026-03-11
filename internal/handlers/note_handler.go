package handlers

import (
	"net/http"

	"backend-AI-Knowledge-Assistant/internal/repositories"
	"backend-AI-Knowledge-Assistant/internal/services"
	"backend-AI-Knowledge-Assistant/pkg/response"

	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	NoteService *services.NoteService
	NoteRepo    *repositories.NoteRepository
}

func NewNoteHandler(noteService *services.NoteService, noteRepo *repositories.NoteRepository) *NoteHandler {
	return &NoteHandler{
		NoteService: noteService,
		NoteRepo:    noteRepo,
	}
}

type noteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *NoteHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	var req noteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	note, err := h.NoteService.CreateNote(userID, req.Title, req.Content)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "note created", note)
}

func (h *NoteHandler) GetAll(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	notes, err := h.NoteRepo.FindAll(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "notes fetched", notes)
}

func (h *NoteHandler) GetByID(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	noteID := c.Param("id")

	note, err := h.NoteRepo.FindByID(noteID, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "note not found")
		return
	}

	response.JSON(c, http.StatusOK, "note fetched", note)
}

func (h *NoteHandler) Update(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	noteID := c.Param("id")

	var req noteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	note, err := h.NoteService.UpdateNote(noteID, userID, req.Title, req.Content)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "note updated", note)
}

func (h *NoteHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	noteID := c.Param("id")

	if err := h.NoteRepo.Delete(noteID, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "note deleted", nil)
}

func (h *NoteHandler) Search(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	q := c.Query("q")

	notes, err := h.NoteRepo.Search(userID, q)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "search result", notes)
}

func (h *NoteHandler) GetSummary(c *gin.Context) {
	noteID := c.Param("id")

	summary, err := h.NoteRepo.GetSummary(noteID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "summary not found")
		return
	}

	response.JSON(c, http.StatusOK, "summary fetched", summary)
}

func (h *NoteHandler) GetKeywords(c *gin.Context) {
	noteID := c.Param("id")

	keywords, err := h.NoteRepo.GetKeywords(noteID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "keywords fetched", keywords)
}