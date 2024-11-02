package comment

type Comment struct {
	ID       int    `json:"id"`
	PostID   int    `json:"post_id"`
	AuthorID int    `json:"author_id"`
	Text     string `json:"text"`
}
