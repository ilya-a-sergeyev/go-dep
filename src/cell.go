package godep

import "math/rand"

type Cell struct {
    instruction		Instruction
    energy		int
    tailId		int
    executorId		int
    dir			Direction
}

func (cell *Cell)takeEnergy() int {
    ret := cell.energy
    cell.energy = 0
    return ret
}

func NewCell() *Cell {
    cell := new(Cell)
    cell.instruction.code = Op_None
    cell.energy = INIT_CELL_ENERGY
    cell.tailId = 0
    cell.executorId = 0
    switch rand.Intn(3) {
    case 0:cell.dir=Forward
    case 1:cell.dir=Back
    case 2:cell.dir=Left
    case 3:cell.dir=Right
    }
    return cell
}

func (cell *Cell)Copy(src *Cell) {
    cell.instruction = src.instruction
    cell.energy = src.energy
    cell.tailId = 0
    cell.executorId = 0
}

func (cell *Cell)Clear() {
    cell.energy = INIT_CELL_ENERGY
    cell.instruction.code = Op_None
    cell.tailId = 0
    cell.executorId = 0
}
