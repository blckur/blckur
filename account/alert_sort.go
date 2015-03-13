package account

import (
	"sort"
)

type AlertsSort struct {
	Alerts []*Alert
}

func (a *AlertsSort) Len() (n int) {
	n = len(a.Alerts)
	return
}

func (a *AlertsSort) Less(i int, j int) (x bool) {
	x = a.Alerts[i].Label < a.Alerts[j].Label
	return
}

func (a *AlertsSort) Swap(i int, j int) {
	a.Alerts[i], a.Alerts[j] = a.Alerts[j], a.Alerts[i]
}

func (a *Account) SortAlerts() {
	srt := &AlertsSort{
		Alerts: a.Alerts,
	}
	sort.Sort(srt)
}
