package token

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"
)

var mu sync.RWMutex

func GenToken(d interface{}) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	var s string
	if d == nil {
		bytes := make([]byte, 52)
		if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
			return "", err
		}
		s = base64.URLEncoding.EncodeToString(bytes)
	} else {
		e, _ := json.Marshal(d)
		s = base64.RawURLEncoding.EncodeToString(e)
	}
	b := md5.Sum([]byte(time.Now().Format("2006-01-02 15:04:05")))
	c := md5.Sum([]byte(s))
	a := fmt.Sprintf("%x.%s.%x", b, s, c)
	return a, nil
}
