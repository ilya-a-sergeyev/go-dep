package godep

import (
    "sync"
    "errors"
)

type IntStack struct {
     lock sync.Mutex 
     s []int
}

func NewIntStack() *IntStack {
    return &IntStack {sync.Mutex{}, make([]int,0), }
}

func (s *IntStack) Push(v int) {
    s.lock.Lock()
    defer s.lock.Unlock()

    s.s = append(s.s, v)
}

func (s *IntStack) Pop() (int, error) {
    s.lock.Lock()
    defer s.lock.Unlock()


    l := len(s.s)
    if l == 0 {
        return 0, errors.New("Empty IntStack")
    }

    res := s.s[l-1]
    s.s = s.s[:l-1]
    return res, nil
}

type CoordStack struct {
     lock sync.Mutex 
     s []Coord
}

func NewCoordStack() *CoordStack {
    return &CoordStack {sync.Mutex{}, make([]Coord,0), }
}

func (s *CoordStack) Push(v Coord) {
    s.lock.Lock()
    defer s.lock.Unlock()

    s.s = append(s.s, v)
}

func (s *CoordStack) Pop() (Coord, error) {
    s.lock.Lock()
    defer s.lock.Unlock()


    l := len(s.s)
    if l == 0 {
        return Coord {0,0}, errors.New("Empty CoordStack")
    }

    res := s.s[l-1]
    s.s = s.s[:l-1]
    return res, nil
}

type CellStack struct {
     lock sync.Mutex 
     s []Cell
}

func NewCellStack() *CellStack {
    return &CellStack {sync.Mutex{}, make([]Cell,0), }
}

func (s *CellStack) Push(v Cell) {
    s.lock.Lock()
    defer s.lock.Unlock()

    s.s = append(s.s, v)
}

func (s *CellStack) Pop() (Cell, error) {
    s.lock.Lock()
    defer s.lock.Unlock()


    l := len(s.s)
    if l == 0 {
        return Cell{}, errors.New("Empty CellStack")
    }

    res := s.s[l-1]
    s.s = s.s[:l-1]
    return res, nil
}

func (s *CellStack) IdxFromTop(idx int) (Cell, error) {
    s.lock.Lock()
    defer s.lock.Unlock()


    l := len(s.s)
    if l == 0 {
        return Cell{}, errors.New("Empty CellStack")
    }

    if l-1 < idx {
        return Cell{}, errors.New("Bad Idx CellStack")
    }

    res := s.s[l-1-idx]

    return res, nil
}

