package context

import "github.com/ezaurum/owlbear"

type Notifier interface {
	Subscribe(string, owlbear.NotificationCallback) (chan owlbear.Event, int64)
	Unsubscribe(string, int64)
	Notify(string, interface{})
}

func (a *app) Subscribe(eventName string, callback owlbear.NotificationCallback) (chan owlbear.Event, int64) {
	return a.eventNotifier.Subscribe(eventName, callback)
}

func (a *app) Unsubscribe(eventName string, subscribeID int64) {
	a.eventNotifier.Unsubscribe(eventName, subscribeID)
}

func (a *app) Notify(eventName string, eventData interface{}) {
	a.eventNotifier.Notify(eventName, eventData)
}
