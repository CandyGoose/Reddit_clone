package comment

import (
	"errors"
	"fmt"
	"sync"
)

type commentService struct {
	mu        sync.Mutex
	comments  []Comment
	commentID int
}

func NewCommentService() Service {
	return &commentService{
		commentID: 1,
		comments:  []Comment{},
	}
}

func (s *commentService) AddComment(postID int, comment Comment) (Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	comment.ID = s.commentID
	comment.PostID = postID
	s.commentID++
	s.comments = append(s.comments, comment)

	return comment, nil
}

func (s *commentService) DeleteComment(postID, commentID, userID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	index := -1
	for i, c := range s.comments {
		if c.PostID == postID && c.ID == commentID {
			if c.AuthorID != userID {
				return errors.New("not authorized to delete this comment")
			}
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("comment with ID %d not found", commentID)
	}

	s.comments = append(s.comments[:index], s.comments[index+1:]...)
	return nil
}
