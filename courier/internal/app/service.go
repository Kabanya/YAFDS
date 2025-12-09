package app

type Service interface {
	Hello() string
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) Hello() string {
	return "Hello world from courier process"
}

var _ Service = (*service)(nil)
