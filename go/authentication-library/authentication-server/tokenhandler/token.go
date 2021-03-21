package token

import (
	"github.com/markstanden/token"
)

func NewToken(m map[string]interface{}) string{
	t := token.NewToken(m)
	return t
}

func Decode(t string) (map[string]interface{}, error){
	return token.Decode(t)
}