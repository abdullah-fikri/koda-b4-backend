package models

import (
	"backend/config"
	"backend/lib"
	"context"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(req RegisterRequest) (*User, error) {
	ctx := context.Background()

	hashedPassword := lib.HashPassword(req.Password)

	var userID int64
	err := config.Db.QueryRow(ctx,
		`INSERT INTO users (email, password, role)
         VALUES ($1, $2, 'user')
         RETURNING id`,
		req.Email, hashedPassword,
	).Scan(&userID)
	if err != nil {
		return nil, err
	}

	_, err = config.Db.Exec(ctx,
		`INSERT INTO profile (users_id, username, phone, address)
         VALUES ($1, $2, $3, $4)`,
		userID, req.Username, req.Phone, req.Address,
	)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:    userID,
		Email: req.Email,
		Role:  "user",
	}, nil
}
func Login(email string) (*User, error) {
	ctx := context.Background()

	var user User
	err := config.Db.QueryRow(ctx,
		`SELECT id, email, password, role
		 FROM users
		 WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.Password, &user.Role)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
