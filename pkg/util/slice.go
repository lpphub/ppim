package util

// RemoveDup 移除重复元素
func RemoveDup[T comparable](s []T) []T {
	seen := make(map[T]bool)
	var result []T
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}
