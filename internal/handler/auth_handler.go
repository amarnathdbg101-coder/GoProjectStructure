package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/amarnathdbg101-coder/GoProjectStructure/pkg/middleware"
	"github.com/amarnathdbg101-coder/GoProjectStructure/internal/model"
	"github.com/amarnathdbg101-coder/GoProjectStructure/pkg/validator"
	"github.com/amarnathdbg101-coder/GoProjectStructure/internal/repository"
	"github.com/amarnathdbg101-coder/GoProjectStructure/internal/service"
)

// AuthHandler handles HTTP requests for user authentication and profiles.
type AuthHandler struct {
	service *service.AuthService
}

// NewAuthHandler constructs a new AuthHandler.
func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register processes a user registration request.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.ValidateStruct(&input); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.Register(r.Context(), input)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			Error(w, http.StatusConflict, "this email address is already registered")
			return
		}
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Created(w, "Account registered successfully", res)
}

// Login authenticates user credentials and returns an auth token.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.ValidateStruct(&input); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.Login(r.Context(), input)
	if err != nil {
		Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	Success(w, "Login successful", res)
}

// Me retrieves the profile of the current authenticated user.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.service.GetProfile(r.Context(), claims.UserID)
	if err != nil {
		Error(w, http.StatusNotFound, "user profile not found")
		return
	}

	Success(w, "Profile fetched", user)
}
