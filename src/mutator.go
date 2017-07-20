package godep

import "sync"

type Mutator struct {
        O interface{};
}

var theMutator *Mutator
var onceMutator sync.Once

func GetMutator() *Mutator {
        onceMutator.Do(func() {
                theMutator = &Mutator{}
        })
        return theMutator
}

func (m *Mutator) chanceToBeRemoved() bool {
    return false
}

func (m *Mutator) chanceToConstantChanged() bool {
    return false
}

func (m *Mutator) chanceToBeRedirected() bool {
    return false
}

func (m *Mutator) chanceToNopInserted() bool {
    return false
}

func (m *Mutator) MutateVector(crd *Coord) {
}
