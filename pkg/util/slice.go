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

func SplitSlice[T any](slice []T, chunkSize int) [][]T {
	// 如果分块大小小于等于 0，直接返回原切片
	if chunkSize <= 0 {
		return [][]T{slice}
	}

	chunkCount := (len(slice) + chunkSize - 1) / chunkSize
	chunks := make([][]T, 0, chunkCount)

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}
