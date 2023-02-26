package time

import "time"

// A service that exists only to wrap time utilities for easy test mocking

var _ TimeService = new(Service)

type TimeService interface {
	Now() time.Time
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s Service) Now() time.Time {
	return time.Now()
}
