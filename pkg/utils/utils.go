package utils

// Функция для разбиения входных данных на батчи
func chunkSlice[T any](slice []T, batchSize int) [][]T {
	var chunks [][]T
	for batchSize < len(slice) {
		slice, chunks = slice[batchSize:], append(chunks, slice[0:batchSize:batchSize])
	}
	chunks = append(chunks, slice)
	return chunks
}
