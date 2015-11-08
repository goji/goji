package pattern

import (
	"goji.io/internal"
	"golang.org/x/net/context"
)

/*
Storage is an implementation of variable and path storage. It attempts to be
efficient and expedient for the common case, but can easily be vendored or
replaced entirely if different performance characteristics are required.

The zero value of Storage is empty and is appropriate for immediate use. Storage
is not safe for concurrent use by multiple goroutines.
*/
type Storage struct {
	shadow shadow
}

/*
Set adds a new variable binding to the Storage object.
*/
func (s *Storage) Set(key string, value interface{}) {
	if s.shadow.overflow != nil {
		s.shadow.overflow[key] = value
	} else if s.shadow.length < len(s.shadow.values) {
		s.shadow.values[s.shadow.length].key = key
		s.shadow.values[s.shadow.length].value = value
		s.shadow.length++
	} else {
		s.shadow.overflow = make(map[string]interface{}, s.shadow.length+1)
		for _, kv := range s.shadow.values {
			s.shadow.overflow[kv.key] = kv.value
		}
		s.shadow.overflow[key] = value
	}
}

/*
SetPath overwrites the path used by Goji's PathPrefix optimization. If a path is
not explicitly provided, the empty path is assumed to have been set.
*/
func (s *Storage) SetPath(path string) {
	s.shadow.path = path
}

/*
Bind returns a context.Context that contains a point-in-time snapshot of the
current state of Storage; subsequent modifictaions to Storage will not be
reflected in the returned Context.

The context.Context returned by this function is safe for concurrent use by
multiple goroutines.
*/
func (s *Storage) Bind(ctx context.Context) context.Context {
	ns := new(shadow)
	*ns = s.shadow
	ns.Context = ctx
	if s.shadow.overflow != nil {
		ns.overflow = make(map[string]interface{}, len(s.shadow.overflow))
		for k, v := range s.shadow.overflow {
			ns.overflow[k] = v
		}
	}
	return ns
}

type shadow struct {
	context.Context
	path   string
	length int
	values [5]struct {
		key   string
		value interface{}
	}
	overflow map[string]interface{}
}

func (s shadow) Value(key interface{}) interface{} {
	switch key {
	case AllVariables:
		var vs map[Variable]interface{}
		if vsi := s.Context.Value(key); vsi == nil {
			if s.length == 0 {
				return nil
			}
			vs = make(map[Variable]interface{})
		} else {
			vs = vsi.(map[Variable]interface{})
		}
		if s.overflow != nil {
			for k, v := range s.overflow {
				vs[Variable(k)] = v
			}
		} else {
			for i := 0; i < s.length; i++ {
				vs[Variable(s.values[i].key)] = s.values[i].value
			}
		}
		return vs
	case internal.Path:
		return s.path
	}

	if k, ok := key.(Variable); ok {
		if s.overflow != nil {
			if v, ok := s.overflow[string(k)]; ok {
				return v
			}
		} else {
			for i := 0; i < s.length; i++ {
				if string(k) == s.values[i].key {
					return s.values[i].value
				}
			}
		}
	}

	return s.Context.Value(key)
}
