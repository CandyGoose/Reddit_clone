package comment

type Repository interface {
	AddComment(postID int, comment Comment) (Comment, error)
	DeleteComment(postID, commentID int) error
}

type Service interface {
	AddComment(postID int, comment Comment) (Comment, error)
	DeleteComment(postID, commentID, userID int) error
}
