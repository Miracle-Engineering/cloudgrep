package template

import (
	"sync"
)

var cache templateCache

type templateCache struct {
	templates map[string]*Template
	l         sync.Mutex
}

func (c *templateCache) get(name string) *Template {
	c.l.Lock()
	defer c.l.Unlock()

	t, has := c.templates[name]
	if !has {
		return nil
	}

	return t
}

func (c *templateCache) put(t *Template) {
	if t == nil {
		return
	}

	c.l.Lock()
	defer c.l.Unlock()

	if c.templates == nil {
		c.templates = make(map[string]*Template)
	}

	c.templates[t.name] = t
}
