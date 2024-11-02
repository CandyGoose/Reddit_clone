package post

import "redditclone/internal/comment"

type Post struct {
	ID        int               `json:"id"`
	Title     string            `json:"title,omitempty"`
	URL       string            `json:"url,omitempty"`
	Text      string            `json:"text,omitempty"`
	Category  string            `json:"category"`
	AuthorID  int               `json:"author_id"`
	Comments  []comment.Comment `json:"comments"`
	Upvotes   int               `json:"upvotes"`
	Downvotes int               `json:"downvotes"`
	Voters    map[int]int
}
