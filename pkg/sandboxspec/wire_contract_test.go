package sandboxspec_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/warpbuilds/warpbuild-agent/pkg/sandboxspec/filesystem"
	"github.com/warpbuilds/warpbuild-agent/pkg/sandboxspec/process"
)

var update = flag.Bool("update", false, "rewrite the wire contract golden file")

func canon(fd protoreflect.FileDescriptor) []string {
	var out []string

	svcs := fd.Services()
	for i := 0; i < svcs.Len(); i++ {
		ms := svcs.Get(i).Methods()
		for j := 0; j < ms.Len(); j++ {
			m := ms.Get(j)
			out = append(out, fmt.Sprintf("RPC /%s/%s in=%s out=%s cstream=%v sstream=%v",
				svcs.Get(i).FullName(), m.Name(), m.Input().FullName(), m.Output().FullName(),
				m.IsStreamingClient(), m.IsStreamingServer()))
		}
	}

	var msgs func(protoreflect.MessageDescriptors)
	msgs = func(mds protoreflect.MessageDescriptors) {
		for i := 0; i < mds.Len(); i++ {
			m := mds.Get(i)
			if m.IsMapEntry() {
				continue
			}
			fs := m.Fields()
			for j := 0; j < fs.Len(); j++ {
				f := fs.Get(j)
				kind := f.Kind().String()
				switch {
				case f.IsMap():
					kind = fmt.Sprintf("map<%s,%s>", f.MapKey().Kind(), f.MapValue().Kind())
				case f.Kind() == protoreflect.MessageKind:
					kind = "message:" + string(f.Message().FullName())
				case f.Kind() == protoreflect.EnumKind:
					kind = "enum:" + string(f.Enum().FullName())
				}
				oneof := ""
				if o := f.ContainingOneof(); o != nil && !o.IsSynthetic() {
					oneof = " oneof=" + string(o.Name())
				}
				out = append(out, fmt.Sprintf("FIELD %s.%s num=%d kind=%s repeated=%v optional=%v%s",
					m.FullName(), f.Name(), f.Number(), kind, f.IsList(), f.HasOptionalKeyword(), oneof))
			}
			msgs(m.Messages())
		}
	}
	msgs(fd.Messages())

	var enums func(protoreflect.EnumDescriptors)
	enums = func(eds protoreflect.EnumDescriptors) {
		for i := 0; i < eds.Len(); i++ {
			e := eds.Get(i)
			vs := e.Values()
			for j := 0; j < vs.Len(); j++ {
				out = append(out, fmt.Sprintf("ENUM %s.%s = %d", e.FullName(), vs.Get(j).Name(), vs.Get(j).Number()))
			}
		}
	}
	enums(fd.Enums())
	var nested func(protoreflect.MessageDescriptors)
	nested = func(mds protoreflect.MessageDescriptors) {
		for i := 0; i < mds.Len(); i++ {
			enums(mds.Get(i).Enums())
			nested(mds.Get(i).Messages())
		}
	}
	nested(fd.Messages())

	sort.Strings(out)
	return out
}

// The sandbox wire format is deliberately identical to e2b envd's so their SDKs
// keep working and old clients need no upgrade. Any diff here is a breaking
// change to every deployed client, not a refactor.
func TestWireContractUnchanged(t *testing.T) {
	var lines []string
	lines = append(lines, canon(process.File_process_process_proto)...)
	lines = append(lines, canon(filesystem.File_filesystem_filesystem_proto)...)
	got := strings.Join(lines, "\n") + "\n"

	golden := filepath.Join("testdata", "wire_contract.golden")

	if *update {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Log("golden updated")
		return
	}

	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("read golden (run with -update to create): %v", err)
	}
	if got == string(want) {
		return
	}

	wantLines := strings.Split(strings.TrimSuffix(string(want), "\n"), "\n")
	inWant := map[string]bool{}
	for _, l := range wantLines {
		inWant[l] = true
	}
	inGot := map[string]bool{}
	for _, l := range lines {
		inGot[l] = true
	}
	for _, l := range wantLines {
		if !inGot[l] {
			t.Errorf("wire contract entry REMOVED: %s", l)
		}
	}
	for _, l := range lines {
		if !inWant[l] {
			t.Errorf("wire contract entry ADDED: %s", l)
		}
	}
}
