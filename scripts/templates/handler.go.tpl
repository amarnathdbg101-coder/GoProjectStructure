package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"goProjectStructure/internal/model"
	"goProjectStructure/internal/pkg/pagination"
	"goProjectStructure/internal/pkg/validator"
	"goProjectStructure/internal/repository"

	"github.com/go-chi/chi/v5"
)

type __PASCAL__Handler struct {
	repo *repository.__PASCAL__Repository
}

func New__PASCAL__Handler(repo *repository.__PASCAL__Repository) *__PASCAL__Handler {
	return &__PASCAL__Handler{repo: repo}
}

func (h *__PASCAL__Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.Create__PASCAL__Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	if err := validator.ValidateStruct(req); err != nil {
		Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	item := &model.__PASCAL__{
		Title:       req.Title,
		Description: req.Description,
	}

	if err := h.repo.Create(r.Context(), item); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Created(w, "__PASCAL__ created successfully", item)
}

func (h *__PASCAL__Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid ID parameter")
		return
	}

	item, err := h.repo.FindByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if item == nil {
		Error(w, http.StatusNotFound, "__PASCAL__ not found")
		return
	}
	Success(w, "__PASCAL__ retrieved successfully", item)
}

func (h *__PASCAL__Handler) List(w http.ResponseWriter, r *http.Request) {
	params := pagination.GetParams(r)
	items, total, err := h.repo.FindAll(r.Context(), params)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	pagination.Respond(w, items, total, params)
}

func (h *__PASCAL__Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid ID parameter")
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, "__PASCAL__ deleted successfully", nil)
}
