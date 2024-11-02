package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"redditclone/internal/comment"
	"redditclone/internal/middleware"
	"redditclone/internal/post"
	"redditclone/internal/user"
)

func main() {
	logger := log.New(os.Stdout, "redditclone: ", log.LstdFlags)

	userService := user.NewUserService(logger)
	postService := post.NewPostService(logger)
	commentService := comment.NewCommentService()

	authHandler := user.NewUserHandler(userService, logger)
	postHandler := post.NewPostHandler(postService, userService, logger)
	commentHandler := comment.NewCommentHandler(commentService, logger)

	router := mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()

	api.HandleFunc("/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/login", authHandler.Login).Methods("POST")

	api.HandleFunc("/posts", postHandler.GetAllPosts).Methods("GET")
	api.HandleFunc("/posts", middleware.JWTMiddleware(postHandler.CreatePost)).Methods("POST")
	api.HandleFunc("/posts/{category}", postHandler.GetPostsByCategory).Methods("GET")
	api.HandleFunc("/post/{postID}", postHandler.GetPostDetails).Methods("GET")
	api.HandleFunc("/post/{postID}", middleware.JWTMiddleware(postHandler.DeletePost)).Methods("DELETE")
	api.HandleFunc("/post/{postID}/upvote", middleware.JWTMiddleware(postHandler.UpvotePost)).Methods("GET")
	api.HandleFunc("/post/{postID}/downvote", middleware.JWTMiddleware(postHandler.DownvotePost)).Methods("GET")
	api.HandleFunc("/post/{postID}/unvote", middleware.JWTMiddleware(postHandler.UnvotePost)).Methods("GET")
	api.HandleFunc("/user/{userLogin}", postHandler.GetPostsByUser).Methods("GET")

	api.HandleFunc("/post/{postID}/comment", middleware.JWTMiddleware(commentHandler.AddComment)).Methods("POST")
	api.HandleFunc("/post/{postID}/comment/{commentID}", middleware.JWTMiddleware(commentHandler.DeleteComment)).Methods("DELETE")

	staticFileDirectory := http.Dir("redditclone/static/")
	staticFileHandler := http.StripPrefix("/static/", http.FileServer(staticFileDirectory))
	router.PathPrefix("/static/").Handler(staticFileHandler)

	logger.Printf("Сервер запущен на %s\n", "http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
