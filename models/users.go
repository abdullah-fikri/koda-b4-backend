package models

import (
	"backend/config"
	"backend/lib"
	"context"
)

type ListUserStruct struct {
	ID             int64  `json:"id"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	Username       string `json:"username"`
	Phone          string `json:"phone"`
	Address        string `json:"address"`
	ProfilePicture string `json:"profile_picture"`
}

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Role     string `json:"role"`
}
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password,omitempty"`
	Username string `json:"username" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Address  string `json:"address" binding:"required"`
	Role     string `json:"-"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
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
func UpdateUser(email string, req RegisterRequest) (*User, error) {
	ctx := context.Background()

	if req.Password != "" {
		hashedPassword := lib.HashPassword(req.Password)

		_, err := config.Db.Exec(ctx,
			`UPDATE users SET password=$1 WHERE email=$2`,
			hashedPassword, email,
		)
		if err != nil {
			return nil, err
		}
	}

	_, err := config.Db.Exec(ctx,
		`UPDATE profile 
		 SET username=$1, phone=$2, address=$3
		 WHERE users_id = (SELECT id FROM users WHERE email=$4)`,
		req.Username, req.Phone, req.Address, email,
	)
	if err != nil {
		return nil, err
	}

	var user User
	err = config.Db.QueryRow(ctx,
		`SELECT id, email, role FROM users WHERE email=$1`,
		email,
	).Scan(&user.ID, &user.Email, &user.Role)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserEmailByID(id int64) (string, error) {
	ctx := context.Background()
	var email string

	err := config.Db.QueryRow(ctx, `SELECT email FROM users WHERE id=$1`, id).Scan(&email)
	if err != nil {
		return "", err
	}

	return email, nil
}

func Forgot(email string) (*User, error) {
	ctx := context.Background()

	var user User
	err := config.Db.QueryRow(ctx,
		`SELECT id,email,role FROM users WHERE email=$1`, email).Scan(&user.ID, &user.Email, &user.Role)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateProfilePicture(userID int64, path string) error {
	ctx := context.Background()

	_, err := config.Db.Exec(ctx,
		`UPDATE profile
         SET profile_picture = $1, updated_at = NOW()
         WHERE users_id = $2`,
		path, userID,
	)

	return err
}

func ListUser() ([]ListUserStruct, error) {
	ctx := context.Background()

	query := `
	SELECT 
		u.id,
		u.email,
		u.role,
		p.username,
		p.phone,
		p.address,
		COALESCE(p.profile_picture, '') AS profile_picture
	FROM users u
	JOIN profile p ON p.users_id = u.id;
	`

	rows, err := config.Db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []ListUserStruct

	for rows.Next() {
		var u ListUserStruct
		err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.Role,
			&u.Username,
			&u.Phone,
			&u.Address,
			&u.ProfilePicture,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}
