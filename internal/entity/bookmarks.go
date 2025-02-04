package entity

type Bookmark struct {
	ID         string `json:"id"`
	BusinessID string `json:"business_id"`
	UserID     string `json:"user_id"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type BookmarksList struct {
	Items []Bookmark `json:"items"`
	Count int        `json:"count"`
}
