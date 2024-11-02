package post

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	"redditclone/internal/comment"
)

type postService struct {
	mu     sync.Mutex
	posts  map[int]Post
	nextID int
	logger *log.Logger
}

func NewPostService(logger *log.Logger) Service {
	return &postService{
		posts:  make(map[int]Post),
		nextID: 1,
		logger: logger,
	}
}

func (s *postService) CreatePost(post Post) (Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post.ID = s.nextID
	s.nextID++
	post.Comments = []comment.Comment{}
	post.Voters = make(map[int]int)
	s.posts[post.ID] = post

	s.logger.Printf("Post created: %+v\n", post)
	return post, nil
}

func (s *postService) GetAllPosts() ([]Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	posts := make([]Post, 0, len(s.posts))
	for _, post := range s.posts {
		posts = append(posts, post)
	}

	return posts, nil
}

func (s *postService) GetPostsByCategory(category string) ([]Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	posts := make([]Post, 0, len(s.posts))
	for _, post := range s.posts {
		if post.Category == category {
			posts = append(posts, post)
		}
	}

	return posts, nil
}

func (s *postService) GetPostByID(id int) (Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, exists := s.posts[id]
	if !exists {
		return Post{}, errors.New("post not found")
	}

	return post, nil
}

func (s *postService) DeletePost(postID, userID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, exists := s.posts[postID]
	if !exists {
		return errors.New("post not found")
	}

	if post.AuthorID != userID {
		return errors.New("not authorized to delete this post")
	}

	delete(s.posts, postID)
	s.logger.Printf("Post deleted: %d\n", postID)
	return nil
}

func (s *postService) UpvotePost(postID, userID int) error {
	post, err := s.GetPostByID(postID)
	if err != nil {
		return err
	}

	if currentVote, exists := post.Voters[userID]; exists && currentVote == 1 {
		return errors.New("already upvoted")
	}

	post.Upvotes++
	post.Voters[userID] = 1
	s.posts[postID] = post
	return nil
}

func (s *postService) DownvotePost(postID, userID int) error {
	post, err := s.GetPostByID(postID)
	if err != nil {
		return err
	}

	if currentVote, exists := post.Voters[userID]; exists && currentVote == -1 {
		return errors.New("already downvoted")
	}

	post.Downvotes++
	post.Voters[userID] = -1
	s.posts[postID] = post
	return nil
}

func (s *postService) UnvotePost(postID, userID int) error {
	post, err := s.GetPostByID(postID)
	if err != nil {
		return err
	}

	if _, exists := post.Voters[userID]; !exists {
		return errors.New("no vote to remove")
	}

	delete(post.Voters, userID)
	s.posts[postID] = post
	return nil
}

func (h *Handler) GetPostsByUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Getting posts by user")

	vars := mux.Vars(r)
	userLogin := vars["userLogin"]

	user, err := h.userService.GetUserByUsername(userLogin)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	posts, err := h.postService.GetPostsByUser(user.ID)
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
