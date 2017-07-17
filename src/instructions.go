package godep

import (
	"strconv"
	"bytes"
)

const (
	// emty cell
	Op_None int = iota
	// do nothing
	Op_Nop
	// start new creature from addres
	Op_Start
	// clear all internal storages and flags
	Op_Reset

	// push world cell value to the internal stack
	Op_Push
	// pop cell value from internal stack to the world by address
	Op_Pop
	// compare top stack value and constant
	Op_CmpTop

	// load and store energy value of the cell  to the internal memory
	Op_GetE
	Op_SetE

	// internal nenory registers manipulations
	Op_Set
	Op_Add
	Op_Rnd
	Op_Mov
	Op_AddReg
	Op_Check
	Op_Len
	Op_Cmp
	Op_Next
	Op_Prev

	// jumps
	Op_Begin
	Op_BreakOnErr
	Op_ContinueOnErr
	Op_BreakOnZ
	Op_ContinueOnZ
	Op_BreakOnG
	Op_ContinueOnG
	Op_BreakOnL
	Op_ContinueOnL
	Op_End

	Op_Last
)

type OperandType int

const ( 
	Ot_None = OperandType(iota)
	Ot_Constant
	Ot_ConstantVector  // vector given by argument
	Ot_IntMemory       // internal memory register address
	Ot_IntMemoryPtr    // external memory address, pointer stored in the internal memory
)

type OpOptions struct {
	baseCost		int
	lengthMul		int
	targetOpType	OperandType
	sourceOpType	OperandType
}

type Instruction struct{
	code		int
	arg1		Coord
	arg2		Coord
}

var opOptions	map[int]OpOptions

func InstructionsPrepare() {

	opOptions[Op_None]   =  OpOptions {ENERGY_COST_FAIL,          0, Ot_None,              Ot_None}
	opOptions[Op_Nop]    =  OpOptions {1,                         0, Ot_None,              Ot_None}
	opOptions[Op_Start]  =  OpOptions {200,                       1, Ot_IntMemoryPtr,      Ot_None}
	opOptions[Op_Reset]  =  OpOptions {10,                        0, Ot_None,              Ot_None}

	// вектор в ячейке внутренней памяти, адрес которой задается операндом
	opOptions[Op_Push ]  =  OpOptions {ENERGY_COST_SENSOR,        1, Ot_IntMemoryPtr,      Ot_ConstantVector}
	opOptions[Op_Pop  ]  =  OpOptions {ENERGY_COST_MODIFICATOR,   2, Ot_IntMemoryPtr,      Ot_ConstantVector}

	opOptions[Op_GetE  ] =  OpOptions {ENERGY_COST_SENSOR,        1, Ot_IntMemory,         Ot_None}
	opOptions[Op_SetE  ] =  OpOptions {ENERGY_COST_MODIFICATOR,   1, Ot_IntMemory,         Ot_None}

	// операнд - константа не связанная с внешним миром
	opOptions[Op_CmpTop] =  OpOptions {ENERGY_COST_INTERNAL,      0, Ot_IntMemory,         Ot_None}
	opOptions[Op_Set]    =  OpOptions {ENERGY_COST_INTERNAL,      0, Ot_IntMemory,         Ot_ConstantVector}
	opOptions[Op_Add]    =  OpOptions {ENERGY_COST_INTERNAL,      0, Ot_IntMemory,         Ot_ConstantVector}
	opOptions[Op_Rnd]    =  OpOptions {ENERGY_COST_INTERNAL,      0, Ot_IntMemory,         Ot_ConstantVector}
	opOptions[Op_Mov]    =  OpOptions {ENERGY_COST_INTERNAL,      0, Ot_IntMemory,         Ot_IntMemory}
	opOptions[Op_AddReg] =  OpOptions {ENERGY_COST_INTERNAL,      0, Ot_IntMemory,         Ot_IntMemory}
	opOptions[Op_Len]    =  OpOptions {ENERGY_COST_INTERNAL,      0, Ot_IntMemory,         Ot_None}
	opOptions[Op_Check]  =  OpOptions {ENERGY_COST_INTERNAL,      0, Ot_IntMemory,         Ot_None}
	opOptions[Op_Cmp  ]  =  OpOptions {ENERGY_COST_INTERNAL,      0, Ot_IntMemory,         Ot_ConstantVector}
	opOptions[Op_Next]   =  OpOptions {ENERGY_COST_INTERNAL,      0, Ot_IntMemory,         Ot_ConstantVector}
	opOptions[Op_Prev]   =  OpOptions {ENERGY_COST_INTERNAL,      0, Ot_IntMemory,         Ot_ConstantVector}

	// операнд - константный вектор во внешний мир относительно текущей ячейки
	opOptions[Op_Begin]         =  OpOptions {ENERGY_COST_INTERNAL,       1, Ot_None,   Ot_None}
	opOptions[Op_BreakOnErr]    =  OpOptions {ENERGY_COST_INTERNAL,       1, Ot_None,   Ot_None}
	opOptions[Op_ContinueOnErr] =  OpOptions {ENERGY_COST_INTERNAL,       1, Ot_None,   Ot_None}
	opOptions[Op_BreakOnZ]      =  OpOptions {ENERGY_COST_INTERNAL,       1, Ot_None,   Ot_None}
	opOptions[Op_ContinueOnZ]   =  OpOptions {ENERGY_COST_INTERNAL,       1, Ot_None,   Ot_None}
	opOptions[Op_BreakOnG]      =  OpOptions {ENERGY_COST_INTERNAL,       1, Ot_None,   Ot_None}
	opOptions[Op_ContinueOnG]   =  OpOptions {ENERGY_COST_INTERNAL,       1, Ot_None,   Ot_None}
	opOptions[Op_BreakOnL]      =  OpOptions {ENERGY_COST_INTERNAL,       1, Ot_None,   Ot_None}
	opOptions[Op_ContinueOnL]   =  OpOptions {ENERGY_COST_INTERNAL,       1, Ot_None,   Ot_None}
	opOptions[Op_End]           =  OpOptions {ENERGY_COST_INTERNAL,       1, Ot_None,   Ot_None}
}

func (i *Instruction)Options() OpOptions {
    return opOptions[i.code];
}

func (i *Instruction)Cost(cell *Cell, crt *Creature, condition bool) int {
    opt := i.Options()
    ee := opt.baseCost

    // constant address
    if opt.targetOpType == Ot_ConstantVector && condition {
        vect := cell.instruction.arg1
        ee += opt.lengthMul*vect.Length()
    } else {
        if opt.targetOpType == Ot_IntMemoryPtr && condition {
            vect := crt.internal_memory[cell.instruction.arg1.x]
            ee += opt.lengthMul*vect.Length()
        }
        if opt.sourceOpType == Ot_IntMemoryPtr && condition {
            vect := crt.internal_memory[cell.instruction.arg2.x]
            ee += opt.lengthMul*vect.Length()
        }
    }

    return ee
}


func (i *Instruction) String(crt *Creature, cell *Cell) string {

	var op_name bytes.Buffer

	switch (i.code) {
	case Op_None:  op_name.WriteString("NONE")
	case Op_Nop:   op_name.WriteString("NOP")
	case Op_Reset: op_name.WriteString("RESET")
	case Op_Start:
		op_name.WriteString("START (")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg1.x))
		op_name.WriteString(":")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg1.y))
		op_name.WriteString(")")
	case Op_Push:
		op_name.WriteString("PUSH ")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg1.x))
		op_name.WriteString(":")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg1.y))
	case Op_Pop:
		op_name.WriteString("POP ")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg1.x))
		op_name.WriteString(":")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg1.y))
	case Op_GetE:
		op_name.WriteString("GETE (")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg1.x))
		op_name.WriteString(") <- ")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg2.x))
		op_name.WriteString(":")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg2.y))
	case Op_SetE:
		op_name.WriteString("SETE (")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg1.x))
		op_name.WriteString(":")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg1.y))
		op_name.WriteString(") <- ")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg2.x))
	case Op_CmpTop:
		op_name.WriteString("CMPTOP ")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg1.x))
	case Op_Set:
		op_name.WriteString("SET ")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg1.x))
		op_name.WriteString(") <- ")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg2.x))
		op_name.WriteString(":")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg2.y))
	case Op_Add:
		op_name.WriteString("ADD (")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg1.x))
		op_name.WriteString(") += ")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg2.x))
		op_name.WriteString(":")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg2.y))
	case Op_Cmp:
		op_name.WriteString("CMP (")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg1.x))
		op_name.WriteString(") with ")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg2.x))
		op_name.WriteString(":")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg2.y))
	case Op_AddReg:
		op_name.WriteString("ADD (")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg1.x))
		op_name.WriteString(") += (")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg2.x))
		op_name.WriteString(")")
	case Op_Mov:
		op_name.WriteString("MOV (")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg1.x))
		op_name.WriteString(") <- (")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg2.x))
		op_name.WriteString(")")
	case Op_Check:
		op_name.WriteString("CHECK ")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg1.x))
	case Op_Next:
		op_name.WriteString("NEXT ")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg1.x))
	case Op_Prev:
		op_name.WriteString("PREV ")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg1.x))
	case Op_Rnd: 
		op_name.WriteString("RND (")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg1.x))
		op_name.WriteString(") <- ")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg2.x))
		op_name.WriteString(":")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg2.y))
	case Op_Len: 
		op_name.WriteString("LEN (")
		op_name.WriteString(strconv.Itoa(cell.instruction.arg1.x))
		op_name.WriteString(")")
	case Op_Begin: op_name.WriteString("BEGIN")
	case Op_BreakOnErr: op_name.WriteString("BREAKONERR")
	case Op_ContinueOnErr: op_name.WriteString("CONTINUEONERR")
	case Op_BreakOnZ: op_name.WriteString("BREAKONZ")
	case Op_ContinueOnZ: op_name.WriteString("CONTINUEONZ")
	case Op_BreakOnG: op_name.WriteString("BREAKONG")
	case Op_ContinueOnG: op_name.WriteString("CONTINUEONG")
	case Op_BreakOnL: op_name.WriteString("BREAKONL")
	case Op_ContinueOnL: op_name.WriteString("CONTINUEONL")
	case Op_End: op_name.WriteString("END")
	default: op_name.WriteString(strconv.Itoa(i.code))
	}
	return op_name.String()
}
