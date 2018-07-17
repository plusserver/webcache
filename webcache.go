// dheilema 2018
// web cache

package webcache

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrUpdateInProgress   = errors.New("Another Update is arleady running")
	ErrWriteWithoutUpdate = errors.New("Writing to struct without StartUpdate()")
)

// the struct to store page data
type CachedPage struct {
	sync.RWMutex
	updating     bool
	content      []byte
	rebuild      []byte
	lastModified time.Time
	maxAge       time.Duration
	requests     uint64
	updates      uint64
}

// create a new CachedPage
// maxAge defines how long the content will stay valid
// after an update
func NewCachedPage(maxAge time.Duration) CachedPage {
	c := CachedPage{}
	c.maxAge = maxAge
	return c
}

// returns if the cached content is valid (not aged out)
func (c *CachedPage) Valid() (v bool) {
	c.RLock()
	defer c.RUnlock()
	t1 := c.lastModified.Add(c.maxAge)
	if t1.After(time.Now()) || c.updating {
		v = true
	}
	return
}

// invalidate the cache
func (c *CachedPage) Clear() (v bool) {
	c.Lock()
	defer c.Unlock()
	c.lastModified = time.Now().Add(-c.maxAge)
	return
}

// the struct can be used as io.Writer
func (c *CachedPage) Write(p []byte) (int, error) {
	c.Lock()
	defer c.Unlock()
	if !c.updating {
		return 0, ErrWriteWithoutUpdate
	}
	n := len(p)
	c.rebuild = append(c.rebuild, p...)
	return n, nil
}

// mark the update transaction as "in progress"
func (c *CachedPage) StartUpdate() error {
	c.Lock()
	defer c.Unlock()
	if c.updating {
		return ErrUpdateInProgress
	}
	c.updating = true
	c.rebuild = []byte{}
	return nil
}

// mark the update transaction as "finished"
func (c *CachedPage) EndUpdate() {
	c.Lock()
	defer c.Unlock()
	c.content = c.rebuild
	c.lastModified = time.Now()
	c.updating = false
	c.updates++
}

// get content
func (c *CachedPage) Get() (out []byte) {
	c.RLock()
	out = c.content
	c.RUnlock()
	c.Lock()
	c.requests++
	c.Unlock()
	return
}

// get metrics of requests handled by cache
// and number of updates
func (c *CachedPage) GetStatistics() (requests, updates uint64) {
	c.RLock()
	defer c.RUnlock()
	requests = c.requests
	updates = c.updates
	return
}

// reset the counter
func (c *CachedPage) ClearStatistics() (requests, updates uint64) {
	c.Lock()
	defer c.Unlock()
	c.requests = 0
	c.updates = 0
	return
}
