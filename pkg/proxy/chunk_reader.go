package proxy

import (
	"bytes"
	"io"
	"sort"
)

func NewUnorderedChunkReader(chunksMap map[int64]ChunkData) (io.Reader, int64) {
	// Sort the chunks by start offset
	offsets := make([]int64, 0, len(chunksMap))
	for offset := range chunksMap {
		offsets = append(offsets, offset)
	}
	sort.Slice(offsets, func(i, j int) bool {
		return chunksMap[offsets[i]].StartOffset < chunksMap[offsets[j]].StartOffset
	})

	// Create readers for each chunk in order
	readers := make([]io.Reader, len(offsets))
	var totalSize int64
	for i, offset := range offsets {
		chunk := chunksMap[offset]
		readers[i] = bytes.NewReader(chunk.Content)
		totalSize += int64(len(chunk.Content))
	}

	// Chain them together with stdlib MultiReader
	return io.MultiReader(readers...), totalSize
}
