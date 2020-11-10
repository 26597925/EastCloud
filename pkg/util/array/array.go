package array

import "sort"

func In(arr []string, val string) bool {
	l := len(arr)
	rel := map[string]int{}
	if l > 100 {
		sort.Strings(arr)
		index := sort.SearchStrings(arr, val)
		if index == l || arr[index] != val {
			return false
		} else {
			return true
		}
	} else {
		for _, s := range arr {
			rel[s] = 1
		}
	}

	_, ok := rel[val]
	return ok
}