package account

import (
	"sort"
)

type alertsSort struct {
	Alerts []*Alert
}

func (a *alertsSort) Len() (n int) {
	n = len(a.Alerts)
	return
}

func (a *alertsSort) Less(i int, j int) (x bool) {
	x = a.Alerts[i].Label < a.Alerts[j].Label
	return
}

func (a *alertsSort) Swap(i int, j int) {
	a.Alerts[i], a.Alerts[j] = a.Alerts[j], a.Alerts[i]
}

func (a *Account) sortAlerts() {
	srt := &alertsSort{
		Alerts: a.Alerts,
	}
	sort.Sort(srt)
}
