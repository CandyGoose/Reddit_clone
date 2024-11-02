package post

type Service interface {
	CreatePost(post Post) (Post, error)
	GetAllPosts() ([]Post, error)
	GetPostsByCategory(category string) ([]Post, error)
	GetPostByID(id int) (Post, error)
	DeletePost(postID, userID int) error
	UpvotePost(postID, userID int) error
	DownvotePost(postID, userID int) error
	UnvotePost(postID, userID int) error
	GetPostsByUser(userID int) ([]Post, error)
}
