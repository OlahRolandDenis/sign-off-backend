package auth

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
)

type SessionStore struct {
	tokens map[string]int64
	mu     sync.RWMutex
}

var Store = &SessionStore{
	tokens: make(map[string]int64),
}

func (Store *SessionStore) CreateToken(agencyID int64) string {
	b := make([]byte, 32)
	rand.Read(b)
	token := base64.URLEncoding.EncodeToString(b)

	Store.mu.Lock()
	defer Store.mu.Unlock()

	Store.tokens[token] = agencyID
	return token
}

func (Store *SessionStore) ValidateToken(token string) (int64, bool) {
	Store.mu.RLock()
	defer Store.mu.RUnlock()

	agencyID, ok := Store.tokens[token]

	return agencyID, ok
}

func (Store *SessionStore) DeleteToken(token string) {
	Store.mu.Lock()
	defer Store.mu.Unlock()

	delete(Store.tokens, token)
}
