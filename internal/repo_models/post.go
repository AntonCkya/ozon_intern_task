package repo_models

type Post struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	UserID      int    `json:"userId"`
	Commentable bool   `json:"commentable"`
}
