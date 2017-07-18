package godep

type World struct {
    size	int
    soup	[]Cell
    creatures	[]Creature
    
}

var (
    TheWorld = World { size:WORLD_SIZE, soup:make([]Cell, WORLD_SIZE*WORLD_SIZE) }
)

type ModelCell struct {
	dir	Direction
	ins	Instruction
}

type Model []ModelCell

func (world *World)GetCell(idx int) *Cell {
    return &world.soup[idx]
}

func (world *World)cellIdx(x, y int) int {
    return x%world.size+(y%world.size)*world.size;
}

func (world *World)GetCellByCoord(c Coord) *Cell {
    return &world.soup[world.cellIdx(c.x, c.y)]
}
