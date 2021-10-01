package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

const codeSep = "\n"
const timeFormat = "2006-01-02 15:04:05.000000"

type login struct {
	authCode  string
	idToken   string
	authTime  time.Time
	tokenTime time.Time
}

type token struct {
	Subject   string
	TimeStamp string
}

func (l *login) isValidCode(time string) bool {
	return l.authTime.Format(timeFormat) == time
}

func (l *login) isNotExpiredCode() bool {
	return time.Now().Sub(l.authTime) <=
		time.Second*time.Duration(config().CodeValiditySecs)
}

type session struct {
	lookup map[string]*login
}

func (s *session) init() {
	s.lookup = make(map[string]*login)
}

func (s *session) authenticate(user string) string {
	t := time.Now()
	ac := fmt.Sprintf("%s%s%s", user, codeSep, t.Format(timeFormat))
	ac = b64.URLEncoding.EncodeToString([]byte(ac))
	s.lookup[user] = &login{
		authCode: ac,
		authTime: t,
	}
	return ac
}

func (s *session) tryGetToken(code string) (bool, string) {
	ok, tkn := false, ""
	ac, err := b64.URLEncoding.DecodeString(code)
	if err != nil {
		log.Printf("Erorr decoding authcode = %s\n", code)
	} else {
		ut := bytes.Split(ac, []byte(codeSep))
		if len(ut) == 2 {
			u := string(ut[0])
			t := string(ut[1])
			l, x := s.lookup[u]
			if x && l != nil {
				// TO DO - Checking existing IdToken
				// Current implmentation always gives New token!
				if l.isValidCode(t) {
					// matches the original session created
					if l.isNotExpiredCode() {
						ok = true
						// create session id token
						tt := time.Now()
						ts := token{
							Subject:   u,
							TimeStamp: tt.Format(timeFormat),
						}
						js, err := json.Marshal(ts)
						if err != nil {
							log.Printf("Error creating token for authcode = %s\n", code)
						}
						tkn = b64.URLEncoding.EncodeToString(js)
						log.Printf("**Token= %s\n**", tkn)
						l.tokenTime = tt
						l.idToken = tkn
					} else {
						log.Printf("Authcode expired = %s\n", code)
					}
				} else {
					log.Printf("Invalid authcode = %s\n", code)
				}
			} else {
				log.Printf("No login session for authcode = %s\n", code)
			}
		} else {
			log.Printf("Invalid authcode format = %s\n", code)
		}
	}
	return ok, tkn
}

/* TEST cURL command to get token
curl -X POST http://localhost:8585/api/token \
 -H 'Content-Type: application/json' \
 -d '{"Code":"QWxhbgoyMDIxLTEwLTAxIDE2OjU3OjI5LjQ1NjY5Mw=="}'
>> eyJTdWJqZWN0IjoiQWxhbiIsIlRpbWVTdGFtcCI6IjIwMjEtMTAtMDEgMTY6NTc6NTguNjgzMjgxIn0=

 => Decode Response
   echo 'eyJTdWJqZWN0IjoiQWxhbiIsIlRpbWVTdGFtcCI6IjIwMjEtMTAtMDEgMTY6NTc6NTguNjgzMjgxIn0=' | base64 -d
 >> {"Subject":"Alan","TimeStamp":"2021-10-01 16:57:58.683281"}
*/
