package repo_models

type Comment struct {
	ID       int    `json:"id"`
	Content  string `json:"content"`
	UserID   int    `json:"userId"`
	PostID   int    `json:"postId"`
	ParentID *int   `json:"parentId,omitempty"`
}
