package intr

import (
	"net/url"

	"github.com/stathat/cmd/stathat/net"
)

func AddComment(id, text string) (string, error) {
	params := url.Values{
		"text": {text},
		"id":   {id},
	}
	var x struct {
		Message string
	}
	if err := net.DefaultAPI.Post("comments", params, &x); err != nil {
		return "", err
	}
	return x.Message, nil
}
