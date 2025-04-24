package services

import "context"

type IService interface {
	Init(provider *ServiceProvider) error
	Run(ctx context.Context) error
}
