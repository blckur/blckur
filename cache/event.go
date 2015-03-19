package cache

const (
	MESSAGE = 1
	RESHARD = 2
)

type Event struct {
	Type int
	Data string
}
