// Init system with before and after constraints.
package requires

import (
	"github.com/dropbox/godropbox/container/set"
	"container/list"
)

var (
	modules = list.New()
)

type Module struct {
	name string
	before set.Set
	after set.Set
	Handler func()
}

func (m *Module) Before(name string) {
	m.before.Add(name)
}

func (m *Module) After(name string) {
	m.after.Add(name)
}

func New(name string) (module *Module) {
	module = &Module{
		name: name,
		before: set.NewSet(),
		after: set.NewSet(),
	}
	modules.PushBack(module)
	return
}

func Init() {
	Loop:
	for {
		i := modules.Front()
		for i != nil {
			module := i.Value.(*Module)

			j := i.Prev()
			for j != nil {
				if module.before.Contains(j.Value.(*Module).name) {
					modules.MoveBefore(i, j)
					continue Loop
				}
				j = j.Prev()
			}

			j = i.Next()
			for j != nil {
				if module.after.Contains(j.Value.(*Module).name) {
					modules.MoveAfter(i, j)
					continue Loop
				}
				j = j.Next()
			}

			i = i.Next()
		}

		break
	}

	i := modules.Front()
	for i != nil {
		i.Value.(*Module).Handler()
		i = i.Next()
	}
}
