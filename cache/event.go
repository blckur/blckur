package cache

const (
	evtMessage = 1
	reshard = 2
)

type event struct {
	Type int
	Data string
}
