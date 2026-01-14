package pubsub

import "strings"

// RPattern builds a RabbitMQ routing pattern
// by concatenating strings with '.'
func RPattern(parts ...string) string {
	return strings.Join(parts, ".")
}
