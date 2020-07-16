package viewhelpers

import (
	"sort"
	"strings"

	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// SortSummarySection sorts the summary sections in name order
func SortSummarySection(sections []component.SummarySection) {
	sort.Slice(sections, func(i, j int) bool {
		s1 := sections[i]
		s2 := sections[j]
		return strings.Compare(s1.Header, s2.Header) < 0
	})
}
