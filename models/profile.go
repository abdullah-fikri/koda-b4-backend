package models

type Profile struct {
	ID        int64  `json:"id"`
	UsersID   int64  `json:"users_id"`
	Username  string `json:"username"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
