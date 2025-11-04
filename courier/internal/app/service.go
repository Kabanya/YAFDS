package app

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Hello() string {
	return "Hello world from courier process"
}
