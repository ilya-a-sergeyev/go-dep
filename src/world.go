package godep

import "sync"

type World struct {
    size	int
    soup	[]Cell
    creatures	[]Creature
    O interface{};
    
}

var theWorld *World
var onceWorld sync.Once

func GetWorld() *World {
        onceWorld.Do(func() {
                theWorld = &World { size:WORLD_SIZE, soup:make([]Cell, WORLD_SIZE*WORLD_SIZE) }
        })
        return theWorld
}

type ModelCell struct {
	dir	Direction
	ins	Instruction
}

type Model []ModelCell

func (world *World)GetCell(idx int) *Cell {
    return &world.soup[idx]
}

func (world *World)CellIdx(x, y int) int {
    return x%world.size+(y%world.size)*world.size;
}

func (world *World)GetCellByCoord(c Coord) *Cell {
    return &world.soup[world.CellIdx(c.x, c.y)]
}
