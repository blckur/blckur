package cache

const (
	MESSAGE = 1
	RESHARD = 2
)

type event struct {
	Type int
	Data string
}
