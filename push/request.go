package push

import (
	"fmt"
	"net/http"
)

const (
	token requestType = iota
	push
	callback
)

func (r requestType) String() string {
	return fmt.Sprintf("mno=%s:: group=%s:: request=%s:", r.MNO(), r.Group(), r.Name())
}

func (r requestType) Endpoint() string {
	panic("implement me")
}

func (r requestType) Method() string {
	switch r {
	default:
		return http.MethodPost
	}
}

func (r requestType) Name() string {
	switch r {
	case token:
		return "token"

	case push:
		return "push"

	case callback:
		return "callback"

	default:
		return ""
	}
}

func (r requestType) MNO() string {
	return "tigo"
}

func (r requestType) Group() string {
	switch r {
	case token:
		return "authorization"

	case push:
		return "collection"

	case callback:
		return "callback"
	}

	return ""
}
