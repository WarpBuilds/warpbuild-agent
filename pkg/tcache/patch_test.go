package tcache

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// literal is a string the way .NET stores it in a dll: a length byte, then UTF-16.
func literal(s string) []byte {
	b := utf16le(s)
	return append([]byte{byte(len(b) + 1)}, b...)
}

func writeDll(t *testing.T, content []byte) (string, string) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bin", "Runner.Worker.dll")
	os.MkdirAll(filepath.Dir(path), 0755)
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}
	return dir, path
}

func TestPatchRunner(t *testing.T) {
	stock := bytes.Join([][]byte{[]byte("head"), literal("ACTIONS_RESULTS_URL"), []byte("mid"), literal("ACTIONS_CACHE_URL"), []byte("tail")}, nil)
	want := bytes.Join([][]byte{[]byte("head"), literal("ZCTIONS_RESULTS_URL"), []byte("mid"), literal("ZCTIONS_CACHE_URL"), []byte("tail")}, nil)

	dir, path := writeDll(t, stock)
	// Patching twice must work: the agent can restart after the first patch.
	for i := 0; i < 2; i++ {
		if err := PatchRunner(dir); err != nil {
			t.Fatal(err)
		}
		if got, _ := os.ReadFile(path); !bytes.Equal(got, want) {
			t.Fatalf("run %d: file not patched as expected", i)
		}
	}
}

func TestPatchRunnerLeavesRunnerStockOnError(t *testing.T) {
	for name, content := range map[string][]byte{
		"one name missing": literal("ACTIONS_RESULTS_URL"),
		"name twice":       bytes.Join([][]byte{literal("ACTIONS_RESULTS_URL"), literal("ACTIONS_CACHE_URL"), literal("ACTIONS_CACHE_URL")}, nil),
		"not a string":     append(append([]byte{0}, utf16le("ACTIONS_RESULTS_URL")...), literal("ACTIONS_CACHE_URL")...),
	} {
		dir, path := writeDll(t, content)
		if err := PatchRunner(dir); err == nil {
			t.Errorf("%s: expected an error", name)
		}
		if got, _ := os.ReadFile(path); !bytes.Equal(got, content) {
			t.Errorf("%s: file was changed", name)
		}
	}
}
