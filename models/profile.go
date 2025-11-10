package models

import (
	"backend/config"
	"context"
)

type Profile struct {
	ID        int64  `json:"id"`
	UsersID   int64  `json:"users_id"`
	Username  string `json:"username"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func UpdateProfilePicture(id string, path string) error {
	ctx := context.Background()

	_, err := config.Db.Exec(ctx,
		`UPDATE profile
		 SET profile_picture = $1
		 WHERE users_id = $2`,
		path, id,
	)

	return err
}
