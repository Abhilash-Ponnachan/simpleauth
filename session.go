package main

import (
	b64 "encoding/base64"
	"fmt"
	"time"
)

type login struct {
	authCode  string
	idToken   string
	authTime  time.Time
	tokenTime time.Time
}

type session struct {
	lookup map[string]login
}

func (s *session) init() {
	s.lookup = make(map[string]login)
}

func (s *session) createAuthCode(user string) string {
	t := time.Now().Format("2006-01-02 15:04:05.000000")
	ac := fmt.Sprintf("%s;%s", user, t)
	ac = b64.URLEncoding.EncodeToString([]byte(ac))
	return ac
}
