package scene

import (
	"math/rand"
)

type Scene struct {
	Seed int64

	Entities []*Entity
}

func New(seed int64) *Scene {
	rand.Seed(seed)

	return &Scene{
		Seed: seed,
	}
}

func (s *Scene) Add(e *Entity) {
	if e.Parameter == nil {
		e.Parameter = make(map[string]interface{})
	}

	s.Entities = append(s.Entities, e)
}
