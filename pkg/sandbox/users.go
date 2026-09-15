//go:build darwin

package sandbox

import (
	"os"
	"os/user"
	"strconv"
	"sync"
	"syscall"
)

// userCache memoises account lookups. On darwin these are never the pure-Go
// path even with CGO disabled — they go through opendirectoryd and allocate a
// multi-kilobyte buffer per call — and a directory listing would otherwise pay
// two of them per entry while resolving the same two or three ids over and over.
type userCache struct {
	fallback *user.User

	mu     sync.RWMutex
	byName map[string]*user.User
	owners map[uint32]string
	groups map[uint32]string
}

func newUserCache(fallback *user.User) *userCache {
	return &userCache{
		fallback: fallback,
		byName:   make(map[string]*user.User),
		owners:   make(map[uint32]string),
		groups:   make(map[uint32]string),
	}
}

// lookup resolves a username to an account, falling back to the guest user for
// an empty or unknown name. The name selects whose home resolves ~ and relative
// paths; it is not a credential.
func (c *userCache) lookup(name string) *user.User {
	if name == "" {
		return c.fallback
	}

	c.mu.RLock()
	u, ok := c.byName[name]
	c.mu.RUnlock()
	if ok {
		return u
	}

	u, err := user.Lookup(name)
	if err != nil {
		u = c.fallback
	}

	c.mu.Lock()
	c.byName[name] = u
	c.mu.Unlock()

	return u
}

func (c *userCache) memo(m map[uint32]string, id uint32, resolve func(string) (string, bool)) string {
	c.mu.RLock()
	name, ok := m[id]
	c.mu.RUnlock()
	if ok {
		return name
	}

	name = strconv.FormatUint(uint64(id), 10)
	if resolved, ok := resolve(name); ok {
		name = resolved
	}

	c.mu.Lock()
	m[id] = name
	c.mu.Unlock()

	return name
}

// ownerGroup resolves a file's uid and gid to names, falling back to the numbers
// when the guest has no matching account.
func (c *userCache) ownerGroup(fi os.FileInfo) (string, string) {
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return "", ""
	}

	owner := c.memo(c.owners, st.Uid, func(id string) (string, bool) {
		u, err := user.LookupId(id)
		if err != nil {
			return "", false
		}

		return u.Username, true
	})
	group := c.memo(c.groups, st.Gid, func(id string) (string, bool) {
		g, err := user.LookupGroupId(id)
		if err != nil {
			return "", false
		}

		return g.Name, true
	})

	return owner, group
}
