package helpers

import "sort"

// offsetInSlice Return the offset of a matching string in a slice or -1 if not found
func offsetInSlice[T comparable](s T, list []T) int {
	for offset, entry := range list {
		if entry == s {
			return offset
		}
	}
	return -1
}

// ListConsistentMerge Consistent merge of TypeList elements, maintaining entries position within the list
// Workaround to TF Plugin SDK issue https://github.com/hashicorp/terraform-plugin-sdk/issues/477
// Taken from https://github.com/alexissavin/terraform-provider-solidserver/blob/master/solidserver/solidserver-helper.go#L62
func ListConsistentMerge[T comparable](old []T, new []T) []interface{} {
	// Step 1 Build local list of member indexed by their offset
	oldOffset := make(map[int]T, len(old))
	diff := make([]T, 0, len(new))
	res := make([]interface{}, 0, len(new))

	for _, n := range new {
		offset := offsetInSlice(n, old)

		if offset != -1 {
			oldOffset[offset] = n
		} else {
			diff = append(diff, n)
		}
	}

	// Merge sorted entries ordered by their offset with the diff array that contain the new ones
	// Step 2 Sort the index
	keys := make([]int, 0, len(old))
	for k := range oldOffset {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// Step 3 build the result
	for _, k := range keys {
		res = append(res, oldOffset[k])
	}
	for _, v := range diff {
		res = append(res, v)
	}

	return res
}
