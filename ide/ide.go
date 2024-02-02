package ide

import (
	"maps"
	"sync"

	"github.com/andyyu2004/gqlt/memosa"
	"github.com/andyyu2004/gqlt/syn"
)

type IDE struct {
	ctx *memosa.Context
	// naive concurrency control
	// acquire read lock while snapshot is used,
	// acquire write lock to update
	lock sync.RWMutex
}

func New() *IDE {
	ctx := memosa.New()
	memosa.Set[inputQuery](ctx, Input{make(map[string]string)})
	return &IDE{ctx, sync.RWMutex{}}
}

func (ide *IDE) Snapshot() (Snapshot, func()) {
	ide.lock.RLock()
	return Snapshot{ide.ctx}, ide.lock.RUnlock
}

// A snapshot of the current state of the IDE.
// All ide operations are performed on a snapshot.
type Snapshot struct {
	ctx *memosa.Context
}

type Changes []Change

type Change interface {
	Apply(*Input)
}

type SetFileContent struct {
	Path    string
	Content string
}

var _ Change = SetFileContent{}

func (s SetFileContent) Apply(input *Input) {
	input.Sources[s.Path] = s.Content
}

func (ide *IDE) Apply(changes Changes) {
	ide.lock.Lock()
	defer ide.lock.Unlock()

	input := Input{maps.Clone(memosa.Fetch[inputQuery](ide.ctx, memosa.InputKey{}).Sources)}
	for _, change := range changes {
		change.Apply(&input)
	}
	memosa.Set[inputQuery](ide.ctx, input)
}

func (sn *Snapshot) Parse(path string) syn.File {
	return memosa.Fetch[parseQuery](sn.ctx, parseKey{path})
}
