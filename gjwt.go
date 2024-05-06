package gjwt

import (
	"encoding/base64"
	"errors"
	"github.com/myfstd/gjwt/cacheEx"
	"github.com/myfstd/gjwt/token"
	"strings"
	"time"
)

type Item struct {
	Data interface{}
	Exp  time.Duration
}

func New(it *Item) (string, error) {
	token, err := token.GenToken(it.Data)
	if err != nil {
		return "", err
	}
	token = token + "." + it.Exp.String()
	cacheEx.Set(token, it.Exp)
	return token, err
}
func Get(token string) (*Item, error) {
	if _, b := cacheEx.Get(token); b {
		v := strings.Split(token, ".")
		if len(v) != 4 {
			return nil, errors.New("token error")
		}
		d, _ := base64.RawURLEncoding.DecodeString(v[1])
		e, _ := time.ParseDuration(v[3])
		go cacheEx.Refresh(token)
		return &Item{Data: d, Exp: e}, nil
	}
	return nil, errors.New("no data")
}
