package nuxeogoclient

import "os"

func GetChunkCount(maxChunkSize, fileSize int64) int {
	remainder := fileSize % maxChunkSize
	quotient := fileSize / maxChunkSize
	if remainder > 0 {
		return int(quotient) + 1
	}
	return int(quotient)
}

func GetChunkSizeAtIndex(maxChunkSize, fileSize int64, chunkIndex int) int64 {
	startOffset := maxChunkSize * int64(chunkIndex)
	if startOffset+maxChunkSize > fileSize {
		return fileSize % maxChunkSize
	}
	return maxChunkSize
}

func ReadFileChunk(file *os.File, maxChunkSize, fileSize int64, chunkIndex int) ([]byte, error) {
	chunkSize := GetChunkSizeAtIndex(maxChunkSize, fileSize, chunkIndex)
	data := make([]byte, chunkSize)
	offset := maxChunkSize * int64(chunkIndex)
	_, err := file.ReadAt(data, offset)
	if err != nil {
		return nil, err
	}
	return data, nil
}
