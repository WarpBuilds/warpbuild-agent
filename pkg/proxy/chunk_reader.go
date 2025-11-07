package proxy

import (
	"io"
	"sort"
)

// UnorderedChunkReader implements io.Reader for reading chunks in offset order
type UnorderedChunkReader struct {
	chunksMap    map[int64]ChunkData
	offsets      []int64
	currentChunk int
	currentPos   int
}

// NewUnorderedChunkReader creates a new ChunkReader from a map of chunks
func NewUnorderedChunkReader(chunksMap map[int64]ChunkData) (*UnorderedChunkReader, int64) {
	// Only sort the offsets, not the chunk data
	offsets := make([]int64, 0, len(chunksMap))
	for offset := range chunksMap {
		offsets = append(offsets, offset)
	}

	sort.Slice(offsets, func(i, j int) bool {
		return chunksMap[offsets[i]].StartOffset < chunksMap[offsets[j]].StartOffset
	})

	var totalSize int64
	for _, offset := range offsets {
		totalSize += int64(len(chunksMap[offset].Content))
	}

	return &UnorderedChunkReader{
		chunksMap:    chunksMap,
		offsets:      offsets,
		currentChunk: 0,
		currentPos:   0,
	}, totalSize
}

// Read implements the io.Reader interface
func (cr *UnorderedChunkReader) Read(p []byte) (n int, err error) {
	if cr.currentChunk >= len(cr.offsets) {
		return 0, io.EOF
	}

	totalRead := 0

	for totalRead < len(p) && cr.currentChunk < len(cr.offsets) {
		chunk := cr.chunksMap[cr.offsets[cr.currentChunk]]
		remaining := len(chunk.Content) - cr.currentPos

		if remaining == 0 {
			// Move to next chunk
			cr.currentChunk++
			cr.currentPos = 0
			continue
		}

		// Copy as much as we can from current chunk
		toCopy := min(len(p)-totalRead, remaining)

		copy(p[totalRead:], chunk.Content[cr.currentPos:cr.currentPos+toCopy])
		totalRead += toCopy
		cr.currentPos += toCopy
	}

	if totalRead == 0 && cr.currentChunk >= len(cr.offsets) {
		return 0, io.EOF
	}

	return totalRead, nil
}
