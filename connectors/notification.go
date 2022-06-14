package connectors

import "fmt"

type NotificationConnector struct{}

func NewNotificationConnector() INotificationConnector {
	return &NotificationConnector{}
}

type INotificationConnector interface {
	Notification(reachLimitBase []float64) error
}

func (c *NotificationConnector) Notification(reachLimitBase []float64) error {
	fmt.Printf("send notification: %v reach to limit\n", reachLimitBase)
	return nil
}
