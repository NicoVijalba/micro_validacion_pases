package auth

import (
	"crypto/subtle"
	"strings"
)

type UserStore struct {
	users map[string]string
}

func NewUserStore(users map[string]string) *UserStore {
	cloned := make(map[string]string, len(users))
	for k, v := range users {
		cloned[strings.TrimSpace(k)] = v
	}
	return &UserStore{users: cloned}
}

func (s *UserStore) Validate(username, password string) bool {
	expected, ok := s.users[strings.TrimSpace(username)]
	if !ok {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(expected), []byte(password)) == 1
}
