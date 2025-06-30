package userhandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"3layerarch/models"
)

//type UserService interface {
//	CreateUser(u models.User) error
//	GetUser(id int) (models.User, error)
//}

type Handler struct {
	Service UserService
}

func New(service UserService) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Empty or unreadable body", http.StatusBadRequest)
		return
	}
	var u models.User
	if err := json.Unmarshal(body, &u); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}
	if err := h.Service.CreateUser(u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	u, err := h.Service.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, _ := json.Marshal(u)
	if _, err := w.Write(b); err != nil {
		fmt.Printf("failed to write response: %v\n", err)
	}
}
