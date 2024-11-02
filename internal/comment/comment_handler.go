package comment

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
	logger  *log.Logger
}

func NewCommentHandler(service Service, logger *log.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

type AddCommentRequest struct {
	Text string `json:"text"`
}

func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Adding a new comment")
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["postID"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var req AddCommentRequest
	if decodeErr := json.NewDecoder(r.Body).Decode(&req); decodeErr != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	comment := Comment{
		Text:     req.Text,
		AuthorID: userID,
	}

	createdComment, err := h.service.AddComment(postID, comment)
	if err != nil {
		http.Error(w, "Could not add comment", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(createdComment); err != nil {
		http.Error(w, "Failed to encode comment", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Deleting comment")
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["postID"])
	
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	commentID, err := strconv.Atoi(vars["commentID"])

	if err != nil {
		http.Error(w, "Invalid IDs", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteComment(postID, commentID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
