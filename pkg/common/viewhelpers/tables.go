package viewhelpers

import (
	"sort"

	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// InitTableFilters sorts the values and sets the selected values to the current values
func InitTableFilters(filters []*component.TableFilter) {
	for _, f := range filters {
		sort.Strings(f.Values)
	}
}

// AddFilterValue adds the filter value if its missing
func AddFilterValue(filter *component.TableFilter, value string) {
	for _, v := range filter.Values {
		if v == value {
			return
		}
	}
	filter.Values = append(filter.Values, value)
}
