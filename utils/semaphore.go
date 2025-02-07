package utils

type GoSemaphore struct {
}

func NewGoSemaphore() *GoSemaphore {
	return &GoSemaphore{}
}

func (s *GoSemaphore) Submit(fn func()) {
	go fn()
}
