package handlers

import (
	"net/http"
	"strings"

	"backend-AI-Knowledge-Assistant/internal/services"
	"backend-AI-Knowledge-Assistant/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		response.Error(c, http.StatusBadRequest, "username dan password wajib diisi")
		return
	}

	user, token, err := h.AuthService.Register(req.Username, req.Password)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "register success", gin.H{
		"user":  user,
		"token": token,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		response.Error(c, http.StatusBadRequest, "username dan password wajib diisi")
		return
	}

	user, token, err := h.AuthService.Login(req.Username, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "username atau password salah")
		return
	}

	response.JSON(c, http.StatusOK, "login success", gin.H{
		"user":  user,
		"token": token,
	})
}