package tcache

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf16"
)

// The runner sets these for every step and overwrites anything we set. Renaming them inside the
// runner makes it write to names no one reads, so the value from our job hook survives.
var renames = [][2]string{
	{"ACTIONS_RESULTS_URL", "ZCTIONS_RESULTS_URL"},
	{"ACTIONS_CACHE_URL", "ZCTIONS_CACHE_URL"},
}

func utf16le(s string) []byte {
	var b []byte
	for _, u := range utf16.Encode([]rune(s)) {
		b = append(b, byte(u), byte(u>>8))
	}
	return b
}

// PatchRunner renames the env var names inside bin/Runner.Worker.dll. It writes only when every
// name is found exactly once, so on any error the runner is left stock. A runner that is already
// patched is a success.
func PatchRunner(runnerDir string) error {
	path := filepath.Join(runnerDir, "bin", "Runner.Worker.dll")
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	patched := 0
	for _, r := range renames {
		from, to := utf16le(r[0]), utf16le(r[1])
		if bytes.Count(data, from) == 0 && bytes.Count(data, to) == 1 {
			patched++
			continue
		}
		i := bytes.Index(data, from)
		// The byte before a .NET string literal is its length: 2 bytes per char plus 1.
		if bytes.Count(data, from) != 1 || i == 0 || int(data[i-1]) != len(from)+1 {
			return fmt.Errorf("%s not found exactly once as a string in %s", r[0], path)
		}
		copy(data[i:], to)
	}
	if patched == len(renames) {
		return nil
	}

	tmp := path + ".warpbuild-tmp"
	if err := os.WriteFile(tmp, data, info.Mode()); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
