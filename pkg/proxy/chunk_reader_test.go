package proxy

import (
	"io"
	"testing"
)

func TestUnorderedChunkReader_Sequential(t *testing.T) {
	chunks := map[int64]ChunkData{
		0:  {StartOffset: 0, EndOffset: 9, Content: []byte("0123456789")},
		20: {StartOffset: 20, EndOffset: 29, Content: []byte("KLMNOPQRST")},
		10: {StartOffset: 10, EndOffset: 19, Content: []byte("abcdefghij")},
	}

	reader, totalSize := NewUnorderedChunkReader(chunks)

	if totalSize != 30 {
		t.Errorf("expected totalSize=30, got %d", totalSize)
	}

	buf := make([]byte, 30)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		t.Fatalf("ReadFull failed: %v", err)
	}

	expected := "0123456789abcdefghijKLMNOPQRST"
	if string(buf) != expected {
		t.Errorf("expected %q, got %q", expected, string(buf))
	}
}

func TestUnorderedChunkReader_OutOfOrder(t *testing.T) {
	// Chunks inserted out of order
	chunks := map[int64]ChunkData{
		20: {StartOffset: 20, EndOffset: 29, Content: []byte("KLMNOPQRST")},
		0:  {StartOffset: 0, EndOffset: 9, Content: []byte("0123456789")},
		10: {StartOffset: 10, EndOffset: 19, Content: []byte("abcdefghij")},
	}

	reader, totalSize := NewUnorderedChunkReader(chunks)

	if totalSize != 30 {
		t.Errorf("expected totalSize=30, got %d", totalSize)
	}

	buf := make([]byte, 30)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		t.Fatalf("ReadFull failed: %v", err)
	}

	// Should read in sorted order regardless of insertion order
	expected := "0123456789abcdefghijKLMNOPQRST"
	if string(buf) != expected {
		t.Errorf("expected %q, got %q", expected, string(buf))
	}
}

func TestUnorderedChunkReader_SmallBufferReads(t *testing.T) {
	chunks := map[int64]ChunkData{
		0:  {StartOffset: 0, EndOffset: 9, Content: []byte("0123456789")},
		10: {StartOffset: 10, EndOffset: 19, Content: []byte("abcdefghij")},
	}

	reader, _ := NewUnorderedChunkReader(chunks)

	// Read in small 5-byte chunks
	buf := make([]byte, 5)

	// First read: "01234"
	n, err := reader.Read(buf)
	if err != nil || n != 5 || string(buf) != "01234" {
		t.Errorf("first read: got n=%d, err=%v, data=%q", n, err, string(buf))
	}

	// Second read: "56789"
	n, err = reader.Read(buf)
	if err != nil || n != 5 || string(buf) != "56789" {
		t.Errorf("second read: got n=%d, err=%v, data=%q", n, err, string(buf))
	}

	// Third read: "abcde" (crosses chunk boundary)
	n, err = reader.Read(buf)
	if err != nil || n != 5 || string(buf) != "abcde" {
		t.Errorf("third read: got n=%d, err=%v, data=%q", n, err, string(buf))
	}

	// Fourth read: "fghij"
	n, err = reader.Read(buf)
	if err != nil || n != 5 || string(buf) != "fghij" {
		t.Errorf("fourth read: got n=%d, err=%v, data=%q", n, err, string(buf))
	}

	// Fifth read: EOF
	n, err = reader.Read(buf)
	if err != io.EOF {
		t.Errorf("expected EOF, got n=%d, err=%v", n, err)
	}
}

func TestUnorderedChunkReader_SingleChunk(t *testing.T) {
	chunks := map[int64]ChunkData{
		0: {StartOffset: 0, EndOffset: 4, Content: []byte("hello")},
	}

	reader, totalSize := NewUnorderedChunkReader(chunks)

	if totalSize != 5 {
		t.Errorf("expected totalSize=5, got %d", totalSize)
	}

	buf := make([]byte, 10)
	n, err := reader.Read(buf)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 5 {
		t.Errorf("expected to read 5 bytes, got %d", n)
	}
	if string(buf[:n]) != "hello" {
		t.Errorf("expected 'hello', got %q", string(buf[:n]))
	}

	// Next read should return EOF
	n, err = reader.Read(buf)
	if err != io.EOF {
		t.Errorf("expected EOF, got n=%d, err=%v", n, err)
	}
}

func TestUnorderedChunkReader_EmptyChunks(t *testing.T) {
	chunks := map[int64]ChunkData{}

	reader, totalSize := NewUnorderedChunkReader(chunks)

	if totalSize != 0 {
		t.Errorf("expected totalSize=0, got %d", totalSize)
	}

	buf := make([]byte, 10)
	n, err := reader.Read(buf)
	if err != io.EOF {
		t.Errorf("expected EOF for empty chunks, got n=%d, err=%v", n, err)
	}
	if n != 0 {
		t.Errorf("expected to read 0 bytes, got %d", n)
	}
}

func TestUnorderedChunkReader_LargeBuffer(t *testing.T) {
	chunks := map[int64]ChunkData{
		0: {StartOffset: 0, EndOffset: 4, Content: []byte("hello")},
	}

	reader, _ := NewUnorderedChunkReader(chunks)

	// Buffer larger than all chunks
	buf := make([]byte, 1000)
	n, err := reader.Read(buf)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 5 {
		t.Errorf("expected to read 5 bytes, got %d", n)
	}
	if string(buf[:n]) != "hello" {
		t.Errorf("expected 'hello', got %q", string(buf[:n]))
	}
}

func TestUnorderedChunkReader_MultipleReadsAcrossChunks(t *testing.T) {
	// Create chunks with varying sizes
	chunks := map[int64]ChunkData{
		11: {StartOffset: 11, EndOffset: 15, Content: []byte("lmnop")},
		0:  {StartOffset: 0, EndOffset: 2, Content: []byte("abc")},
		8:  {StartOffset: 8, EndOffset: 9, Content: []byte("ij")},
		3:  {StartOffset: 3, EndOffset: 7, Content: []byte("defgh")},
		10: {StartOffset: 10, EndOffset: 10, Content: []byte("k")},
	}

	reader, totalSize := NewUnorderedChunkReader(chunks)

	if totalSize != 16 {
		t.Errorf("expected totalSize=16, got %d", totalSize)
	}

	// Read all data
	result, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}

	expected := "abcdefghijklmnop"
	if string(result) != expected {
		t.Errorf("expected %q, got %q", expected, string(result))
	}
}

func TestUnorderedChunkReader_PartialRead(t *testing.T) {
	chunks := map[int64]ChunkData{
		0: {StartOffset: 0, EndOffset: 9, Content: []byte("0123456789")},
	}

	reader, _ := NewUnorderedChunkReader(chunks)

	// Read only part of the data
	buf := make([]byte, 5)
	n, err := reader.Read(buf)
	if err != nil || n != 5 {
		t.Errorf("first read: got n=%d, err=%v", n, err)
	}
	if string(buf) != "01234" {
		t.Errorf("expected '01234', got %q", string(buf))
	}

	// Read rest
	buf2 := make([]byte, 10)
	n, err = reader.Read(buf2)
	if err != nil || n != 5 {
		t.Errorf("second read: got n=%d, err=%v", n, err)
	}
	if string(buf2[:n]) != "56789" {
		t.Errorf("expected '56789', got %q", string(buf2[:n]))
	}

	// EOF
	n, err = reader.Read(buf2)
	if err != io.EOF {
		t.Errorf("expected EOF, got n=%d, err=%v", n, err)
	}
}

func TestUnorderedChunkReader_RealWorldScenario(t *testing.T) {
	// Simulate real BuildKit chunks (32MB chunks)
	const chunkSize = 32 * 1024 * 1024

	// Create 3 chunks out of order
	chunk1 := make([]byte, chunkSize)
	for i := range chunk1 {
		chunk1[i] = 'A'
	}

	chunk2 := make([]byte, chunkSize)
	for i := range chunk2 {
		chunk2[i] = 'B'
	}

	chunk3 := make([]byte, 1024*1024) // Last chunk smaller
	for i := range chunk3 {
		chunk3[i] = 'C'
	}

	chunks := map[int64]ChunkData{
		chunkSize * 2: {StartOffset: chunkSize * 2, EndOffset: chunkSize*2 + int64(len(chunk3)) - 1, Content: chunk3},
		0:             {StartOffset: 0, EndOffset: chunkSize - 1, Content: chunk1},
		chunkSize:     {StartOffset: chunkSize, EndOffset: chunkSize*2 - 1, Content: chunk2},
	}

	reader, totalSize := NewUnorderedChunkReader(chunks)

	expectedSize := int64(len(chunk1) + len(chunk2) + len(chunk3))
	if totalSize != expectedSize {
		t.Errorf("expected totalSize=%d, got %d", expectedSize, totalSize)
	}

	// Read first chunk
	buf1 := make([]byte, chunkSize)
	n, err := io.ReadFull(reader, buf1)
	if err != nil || n != chunkSize {
		t.Fatalf("failed to read first chunk: n=%d, err=%v", n, err)
	}
	if buf1[0] != 'A' || buf1[chunkSize-1] != 'A' {
		t.Error("first chunk content mismatch")
	}

	// Read second chunk
	buf2 := make([]byte, chunkSize)
	n, err = io.ReadFull(reader, buf2)
	if err != nil || n != chunkSize {
		t.Fatalf("failed to read second chunk: n=%d, err=%v", n, err)
	}
	if buf2[0] != 'B' || buf2[chunkSize-1] != 'B' {
		t.Error("second chunk content mismatch")
	}

	// Read third chunk
	buf3 := make([]byte, 1024*1024)
	n, err = io.ReadFull(reader, buf3)
	if err != nil || n != 1024*1024 {
		t.Fatalf("failed to read third chunk: n=%d, err=%v", n, err)
	}
	if buf3[0] != 'C' || buf3[len(buf3)-1] != 'C' {
		t.Error("third chunk content mismatch")
	}

	// Should be EOF now
	_, err = reader.Read(make([]byte, 1))
	if err != io.EOF {
		t.Errorf("expected EOF, got %v", err)
	}
}
