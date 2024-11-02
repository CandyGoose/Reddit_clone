package post

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"redditclone/internal/user"
)

type Handler struct {
	postService Service
	userService user.Service
	logger      *log.Logger
}

func NewPostHandler(postService Service, userService user.Service, logger *log.Logger) *Handler {
	return &Handler{
		postService: postService,
		userService: userService,
		logger:      logger,
	}
}

type CreatePostRequest struct {
	Title    string `json:"title,omitempty"`
	URL      string `json:"url,omitempty"`
	Text     string `json:"text,omitempty"`
	Category string `json:"category"`
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Creating a new post")

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	post := Post{
		Title:    req.Title,
		URL:      req.URL,
		Text:     req.Text,
		Category: req.Category,
		AuthorID: userID,
	}

	createdPost, err := h.postService.CreatePost(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(createdPost); err != nil {
		http.Error(w, "Failed to encode post", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Getting all posts")

	posts, err := h.postService.GetAllPosts()
	if err != nil {
		http.Error(w, "Could not retrieve posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, "Failed to encode posts", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetPostsByCategory(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Getting posts by category")

	vars := mux.Vars(r)
	category := vars["category"]

	posts, err := h.postService.GetPostsByCategory(category)
	if err != nil {
		http.Error(w, "Could not retrieve posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, "Failed to encode posts", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetPostDetails(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Getting post details")

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["postID"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.postService.GetPostByID(postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, "Failed to encode post", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Deleting post")

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

	err = h.postService.DeletePost(postID, userID)
	if err != nil {
		if err.Error() == "not authorized to delete this post" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpvotePost(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Upvoting a post")

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

	err = h.postService.UpvotePost(postID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DownvotePost(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Downvoting a post")

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

	err = h.postService.DownvotePost(postID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UnvotePost(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Unvoting a post")

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

	err = h.postService.UnvotePost(postID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *postService) GetPostsByUser(userID int) ([]Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	userPosts := []Post{}
	for _, post := range s.posts {
		if post.AuthorID == userID {
			userPosts = append(userPosts, post)
		}
	}

	return userPosts, nil
}
