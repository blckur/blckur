package gdefer

var (
	handlers []func()
)

func Defer(handler func()) {
	handlers = append(handlers, handler)
}

func End() {
	for _, handler := range handlers {
		handler()
	}
}

func init() {
	handlers = []func(){}
}
