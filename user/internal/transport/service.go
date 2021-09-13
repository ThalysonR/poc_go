package transport

type Service interface {
	Start() error
	Stop() error
}
