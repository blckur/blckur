package account

import (
	"sort"
)

type filtersSort struct {
	Filters []*Filter
}

func (f *filtersSort) Len() (n int) {
	n = len(f.Filters)
	return
}

func (f *filtersSort) Less(i int, j int) (x bool) {
	x = f.Filters[i].Label < f.Filters[j].Label
	return
}

func (f *filtersSort) Swap(i int, j int) {
	f.Filters[i], f.Filters[j] = f.Filters[j], f.Filters[i]
}

func (f *Account) sortFilters() {
	srt := &filtersSort{
		Filters: f.Filters,
	}
	sort.Sort(srt)
}
