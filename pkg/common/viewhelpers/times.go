package viewhelpers

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ToTimeMarkdown returns a formatted time as markdown
func ToTimeMarkdown(mt *metav1.Time) string {
	if mt == nil {
		return ""
	}
	return mt.Time.Format(time.RFC822)
}
