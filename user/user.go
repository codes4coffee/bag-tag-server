package user

import "github.com/google/uuid"

type User struct {
	Id               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	passcode         string
	sessionToken     string
	tokenGeneratedAt int64
}

func userFinder(users []User, requestedUser string) int {
	for idx, user := range users {
		if user.Name == requestedUser {
			return idx
		}
	}
	return -1
}

func findUserBySessionToken(token string, users []User) *User {
	for _, user := range users {
		if user.sessionToken == token {
			return &user
		}
	}
	return nil
}
