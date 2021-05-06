package ussd

import (
	"context"
)

//type SubscriberNameRequest struct {
//	Msisdn              string
//	CompanyName         string
//	CustomerReferenceID string
//}

type NameCheckResponse struct {
	Result    string
	ErrorCode string
	ErrorDesc string
	Msisdn    string
	Flag      string
	Content   string
}

type NameCheckFunc func(context.Context, SubscriberNameRequest) (NameCheckResponse, error)

func (f NameCheckFunc) NameCheck(ctx context.Context, request SubscriberNameRequest) (resp NameCheckResponse, err error) {
	return f(ctx, request)
}

type NameChecker interface {
	NameCheck(ctx context.Context, request SubscriberNameRequest) (resp NameCheckResponse, err error)
}
