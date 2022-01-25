package model

import (
	"sync"

	"github.com/rs/zerolog/log"
)

// SilentStack represents a stacks of components.
type SilentStack struct {
	components []Component
	mx         sync.RWMutex
}

// NewSilentStack returns a new initialized stack.
func NewSilentStack() *SilentStack {
	return &SilentStack{}
}

// Flatten returns a string representation of the component stack.
func (s *SilentStack) Flatten() []string {
	s.mx.RLock()
	defer s.mx.RUnlock()

	ss := make([]string, len(s.components))
	for i, c := range s.components {
		ss[i] = c.Name()
	}
	return ss
}

// Push adds a new item.
func (s *SilentStack) Push(c Component) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.components = append(s.components, c)
}

// Pop removed the top item and returns it.
func (s *SilentStack) Pop() (Component, bool) {
	if s.Empty() {
		return nil, false
	}

	var c Component
	s.mx.Lock()
	{
		c = s.components[len(s.components)-1]
		s.components = s.components[:len(s.components)-1]
	}
	s.mx.Unlock()

	return c, true
}

// Peek returns stack state.
func (s *SilentStack) Peek() []Component {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.components
}

// Clear clears out the stack using pops.
func (s *SilentStack) Clear() {
	for range s.components {
		s.Pop()
	}
}

// Empty returns true if the stack is empty.
func (s *SilentStack) Empty() bool {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return len(s.components) == 0
}

// IsLast indicates if stack only has one item left.
func (s *SilentStack) IsLast() bool {
	return len(s.components) == 1
}

// Previous returns the previous component if any.
func (s *SilentStack) Previous() Component {
	if s.Empty() {
		return nil
	}
	if s.IsLast() {
		return s.Top()
	}

	return s.components[len(s.components)-2]
}

// Top returns the top most item or nil if the stack is empty.
func (s *SilentStack) Top() Component {
	if s.Empty() {
		return nil
	}

	return s.components[len(s.components)-1]
}

// ----------------------------------------------------------------------------
// Helpers...

// Dump prints out the stack.
func (s *SilentStack) Dump() {
	log.Debug().Msgf("--- Stack Dump %p---", s)
	for i, c := range s.components {
		log.Debug().Msgf("%d -- %s -- %#v", i, c.Name(), c)
	}
	log.Debug().Msg("------------------")
}
