package ussd

import (
	"context"
)

type NameCheckRequest struct {
	Msisdn              string
	CompanyName         string
	CustomerReferenceID string
}

type NameCheckResponse struct {
	Result    string
	ErrorCode string
	ErrorDesc string
	Msisdn    string
	Flag      string
	Content   string
}

type NameCheckFunc func(context.Context, NameCheckRequest) (NameCheckResponse, error)

func (f NameCheckFunc) NameCheck(ctx context.Context, request NameCheckRequest) (resp NameCheckResponse, err error) {
	return f(ctx, request)
}

type NameChecker interface {
	NameCheck(ctx context.Context, request NameCheckRequest) (resp NameCheckResponse, err error)
}
