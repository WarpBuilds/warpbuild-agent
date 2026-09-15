//go:build darwin

package sandbox

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"connectrpc.com/connect"
	"github.com/fsnotify/fsnotify"
	"google.golang.org/protobuf/types/known/timestamppb"

	rpc "github.com/warpbuilds/warpbuild-agent/pkg/sandboxspec/filesystem"
	"github.com/warpbuilds/warpbuild-agent/pkg/sandboxspec/filesystem/filesystemconnect"
)

type filesystemService struct {
	filesystemconnect.UnimplementedFilesystemHandler

	users *userCache
}

func newFilesystemService(users *userCache) *filesystemService {
	return &filesystemService{users: users}
}

func (s *filesystemService) entryInfo(path string, fi os.FileInfo) *rpc.EntryInfo {
	e := &rpc.EntryInfo{
		Name:         fi.Name(),
		Path:         path,
		Size:         fi.Size(),
		Mode:         uint32(fi.Mode().Perm()),
		Permissions:  fi.Mode().String(),
		ModifiedTime: timestamppb.New(fi.ModTime()),
		Type:         rpc.FileType_FILE_TYPE_FILE,
	}

	switch {
	case fi.Mode()&os.ModeSymlink != 0:
		e.Type = rpc.FileType_FILE_TYPE_SYMLINK
		if target, err := os.Readlink(path); err == nil {
			e.SymlinkTarget = &target
		}
	case fi.IsDir():
		e.Type = rpc.FileType_FILE_TYPE_DIRECTORY
	}

	e.Owner, e.Group = s.users.ownerGroup(fi)

	return e
}

func (s *filesystemService) Stat(
	ctx context.Context,
	req *connect.Request[rpc.StatRequest],
) (*connect.Response[rpc.StatResponse], error) {
	u := s.users.lookup(basicAuthUsername(req.Header()))
	path := expandPath(req.Msg.GetPath(), u)

	fi, err := os.Lstat(path)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	return connect.NewResponse(&rpc.StatResponse{Entry: s.entryInfo(path, fi)}), nil
}

func (s *filesystemService) ListDir(
	ctx context.Context,
	req *connect.Request[rpc.ListDirRequest],
) (*connect.Response[rpc.ListDirResponse], error) {
	u := s.users.lookup(basicAuthUsername(req.Header()))
	root := expandPath(req.Msg.GetPath(), u)

	depth := int(req.Msg.GetDepth())
	if depth == 0 {
		depth = 1
	}

	resolved, err := filepath.EvalSymlinks(root)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	fi, err := os.Stat(resolved)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	if !fi.IsDir() {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errNotADirectory)
	}

	var entries []*rpc.EntryInfo
	walkErr := filepath.WalkDir(resolved, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			// Entries can vanish mid-walk; skip rather than fail the whole listing.
			if os.IsNotExist(err) {
				return nil
			}

			return err
		}
		if p == resolved {
			return nil
		}

		rel, relErr := filepath.Rel(resolved, p)
		if relErr != nil {
			return nil
		}
		level := strings.Count(rel, string(filepath.Separator)) + 1
		if level > depth {
			if d.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		info, infoErr := d.Info()
		if infoErr != nil {
			return nil
		}
		// Report the path the caller asked about, not the symlink-resolved one.
		entries = append(entries, s.entryInfo(filepath.Join(root, rel), info))

		// Descending further would read and sort a directory whose every child
		// the depth test above will reject.
		if d.IsDir() && level >= depth {
			return filepath.SkipDir
		}

		return nil
	})
	if walkErr != nil {
		return nil, connect.NewError(connect.CodeInternal, walkErr)
	}

	return connect.NewResponse(&rpc.ListDirResponse{Entries: entries}), nil
}

func (s *filesystemService) WatchDir(
	ctx context.Context,
	req *connect.Request[rpc.WatchDirRequest],
	stream *connect.ServerStream[rpc.WatchDirResponse],
) error {
	u := s.users.lookup(basicAuthUsername(req.Header()))
	root := expandPath(req.Msg.GetPath(), u)

	fi, err := os.Stat(root)
	if err != nil {
		return connect.NewError(connect.CodeNotFound, err)
	}
	if !fi.IsDir() {
		return connect.NewError(connect.CodeFailedPrecondition, errNotADirectory)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}
	defer watcher.Close()

	if err := watcher.Add(root); err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}
	if req.Msg.GetRecursive() {
		_ = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
			if err != nil || !d.IsDir() || p == root {
				return nil //nolint:nilerr
			}

			return watcher.Add(p)
		})
	}

	// The first frame always announces the watch is armed, so a client knows
	// when it is safe to start making changes it expects to observe.
	if err := stream.Send(&rpc.WatchDirResponse{
		Event: &rpc.WatchDirResponse_Start{Start: &rpc.WatchDirResponse_StartEvent{}},
	}); err != nil {
		return err
	}

	sendEvent := func(ev fsnotify.Event) error {
		for _, out := range filesystemEvents(root, ev, req.Msg.GetIncludeEntry(), s) {
			if err := stream.Send(&rpc.WatchDirResponse{
				Event: &rpc.WatchDirResponse_Filesystem{Filesystem: out},
			}); err != nil {
				return err
			}
		}

		return nil
	}

	go func() {
		<-ctx.Done()
		watcher.Close()
	}()

	return pumpWithKeepalive(ctx, watcher.Events, keepAliveInterval(req.Header()), sendEvent,
		func() error {
			return stream.Send(&rpc.WatchDirResponse{
				Event: &rpc.WatchDirResponse_Keepalive{Keepalive: &rpc.WatchDirResponse_KeepAlive{}},
			})
		}, nil)
}

// filesystemEvents fans one fsnotify event out into one message per op bit, so a
// write-and-chmod arrives as two events rather than a compound one.
func filesystemEvents(root string, ev fsnotify.Event, includeEntry bool, s *filesystemService) []*rpc.FilesystemEvent {
	name, err := filepath.Rel(root, ev.Name)
	if err != nil {
		name = ev.Name
	}

	ops := []struct {
		op  fsnotify.Op
		typ rpc.EventType
	}{
		{fsnotify.Create, rpc.EventType_EVENT_TYPE_CREATE},
		{fsnotify.Write, rpc.EventType_EVENT_TYPE_WRITE},
		{fsnotify.Remove, rpc.EventType_EVENT_TYPE_REMOVE},
		{fsnotify.Rename, rpc.EventType_EVENT_TYPE_RENAME},
		{fsnotify.Chmod, rpc.EventType_EVENT_TYPE_CHMOD},
	}

	var out []*rpc.FilesystemEvent
	for _, o := range ops {
		if !ev.Op.Has(o.op) {
			continue
		}
		fe := &rpc.FilesystemEvent{Name: name, Type: o.typ}
		// Remove and rename-away leave nothing to stat, and the path may already
		// hold a replacement, so an entry is only attached where it is meaningful.
		if includeEntry && (o.typ == rpc.EventType_EVENT_TYPE_CREATE ||
			o.typ == rpc.EventType_EVENT_TYPE_WRITE ||
			o.typ == rpc.EventType_EVENT_TYPE_CHMOD) {
			if fi, err := os.Lstat(ev.Name); err == nil {
				fe.Entry = s.entryInfo(ev.Name, fi)
			}
		}
		out = append(out, fe)
	}

	return out
}
