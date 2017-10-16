package pubsub

type ReceiveFunc func(issue interface{}, value interface{})

// SubFunc wrap receive func for comparable type.
type SubFunc struct {
	f ReceiveFunc
}

func NewSubFunc(f func(issue interface{}, value interface{})) *SubFunc {
	return &SubFunc{f}
}

func (s *SubFunc) Receive(issue interface{}, value interface{}) {
	s.f(issue, value)
}
