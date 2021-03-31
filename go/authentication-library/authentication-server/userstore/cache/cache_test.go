package cache

import (
	"testing"

	"github.com/markstanden/authentication"
)

func TestFind(t *testing.T) {
	tests := []struct {
		desc  string
		key   string
		value *authentication.User
	}{
		{
			desc: "",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

		})
	}
}
