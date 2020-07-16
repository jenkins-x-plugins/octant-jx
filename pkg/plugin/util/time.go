package util

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DurationString returns the duration between start and end time as string
func DurationString(start *metav1.Time, end *metav1.Time) string {
	if start == nil || end == nil {
		return ""
	}
	return end.Sub(start.Time).Round(time.Second).String()
}
