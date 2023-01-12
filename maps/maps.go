package maps

// MergeMap 将 b 合并到 a 中
// 如果 b 与 a 的键冲突，将覆盖 a
func MergeMap[K comparable, V any](a, b map[K]V) {
	for k, v := range b {
		a[k] = v
	}
}
