package godep

type Fingerprint []Coord

type Creature struct {
	Id				int
	entry				Coord
	ptr				Coord
	flagz, flagg, flagl, flage	bool

// энергия задается при старте родителем и убывает по мере выполения операции
// после ее исчерпания наступает смерть существа
	energy		int
// энергия которую можно передать ребенку при рождении
	child_energy	int

// стек и память виртуальной машины индивидуальны для каждого существа
	internal_stack		CellStack
	internal_memory		[INTERNAL_MEMORY_CAPACITY]Coord
// стек функциональных блоков
	points			CoordStack

// отпечаток жизни существа - для статистики и подчистки территории после его смерти
	lifetime		int
	fingerprint		Fingerprint

// еще не рожденный ребенок, необходимо для очистки этих клеток
// при смерти родителя до старта ребенка
	embrion			Fingerprint

// флаг произошедшей мутации (ограничиваем 1 на проход)
	m_flag			bool
}

var CrIdCnt int

func NewCreature(entryPoint Coord, _eng int) *Creature {

    crt := new(Creature)
    CrIdCnt++
    crt.Id = CrIdCnt
    crt.entry = entryPoint
    crt.ptr = entryPoint
    crt.energy = _eng
    crt.child_energy = 0
    crt.m_flag = false
    crt.lifetime = 0

    //Log::Not << "<<CRT " << Id << " [" << ptr.x << ":" << ptr.y << " E "<< "] was born.>>" << log4cpp::eol;
    //counter1 = 0
    //counter2 = 0
    return crt
}

func (crt *Creature) Destroy() {
    //Log::Not <<
    //    "<<CRT " << Id << " is died (length "<< fingerprint.size() << "). His life was " << lifetime << " instructions long. >>" << log4cpp::eol;

    for i := range crt.fingerprint {
        cell := GetWorld().GetCellByCoord(crt.fingerprint[i])
        //Log::Inf << "<<!!! CRT " << Id << " " << step%World::worldSize << ":"<< step/World::worldSize << ">>" << log4cpp::eol;
        cell.Clear()
    }

    crt.fingerprint = crt.fingerprint[:0]
    //crt.DoAbort()
}

func (crt *Creature)DoAbort() {
    if len(crt.embrion)>0 {
        //Log::Inf << " CRT " << Id << " has " << embrion.size() << " embrional cells." << log4cpp::eol;
        for i := range crt.embrion {
            cell := GetWorld().GetCellByCoord(crt.embrion[i])
            if cell.tailId == crt.Id {
                    //Log::Inf <<
                    //    "<<!!! CRT " << Id << " " << step%World::worldSize << ":"<< step/World::worldSize << ">>" << log4cpp::eol;
                cell.Clear()
            }
        }
        crt.embrion = crt.embrion[:0]
    }
}


func (crt *Creature)JumpToBegin() bool {
    if crt.points.Len()==0 {
        crt.flage = true
        //Log::Inf << "Failed";
        return false
    } 
    nPtr, err := crt.points.Pop()
    if err != nil {
	crt.ptr = nPtr;
    }
    //Log::Inf << "[" << ptr.x << ":" << ptr.y << "]";
    return true
}

func (crt *Creature)jumpToEnd() bool {
    //Log::Inf << log4cpp::eol;
    ins := GetWorld().GetCellByCoord(crt.ptr).instruction
    dir := GetWorld().GetCellByCoord(crt.ptr).dir

    for !(ins.code == Op_End) && !(ins.code == Op_None) {
        //Log::Inf << "CRT " << Id << " [" << ptr.x << ":" << ptr.y << "] " << static_cast<int>(dir);
        //Log::Inf << log4cpp::eol;
        crt.ptr = crt.ptr.Next(dir)
        ins = GetWorld().GetCellByCoord(crt.ptr).instruction
        dir = GetWorld().GetCellByCoord(crt.ptr).dir
    }

    if crt.points.Len()>0 {
        crt.points.Pop()
    }

    return true
}

func (crt *Creature)moveBy(steps int) {
    if steps<0 {
        cnt := -steps;
        for cnt>0 {
            tcell := GetWorld().GetCellByCoord(crt.ptr)
            crt.ptr.Dec(tcell.dir)
            cnt--
        }
        return
    }
    if steps>0 {
        cnt := steps;
        for cnt>0 {
            tcell := GetWorld().GetCellByCoord(crt.ptr)
            crt.ptr.Inc(tcell.dir)
            cnt--
        }
    }
}

func (crt *Creature)setFlags(value int) {
    if value > 0 {
	crt.flagz = false
	crt.flagg = true
	crt.flagl = false
    } else if (value < 0) {
	crt.flagz = false
	crt.flagg = false
	crt.flagl = true
    } else {
	crt.flagz = true
	crt.flagg = false
	crt.flagl = false
    }
}

//
// Place of mutations
//
func (crt *Creature)ApplyPop(targetCoord *Coord, target *Cell, value *Instruction, src_dir Direction) {

    tgt_dir := src_dir

    //  mutation: deletion
    for {

        if GetMutator().chanceToBeRemoved() || crt.m_flag {
            break
        }

        // remove only Nop's
        if value.code != Op_Nop {
            break
        }

        // корректируем указатель на цель, чтобы дальнейшее копирование выполнялось корректно
        cell := GetWorld().GetCellByCoord(crt.ptr)
        vect := &crt.internal_memory[cell.instruction.arg1.x]
        *vect = vect.Prev(src_dir)
        crt.m_flag = true

        //to := crt.ptr.Add(vect.x, vect.y)
        //Log::Not << "CRT " << Id << " <<Op deleted in [" << to.x << ":" << to.y << "]" << log4cpp::eol;
        return
    }

    // mutation: direction shift
    for {
        if GetMutator().chanceToBeRedirected() || crt.m_flag {
            break
        }

        switch src_dir {
        case Forward:tgt_dir = Right
        case Right:tgt_dir = Back
        case Back:tgt_dir = Left
        case Left:tgt_dir = Forward
        }
        //Log::Not << "CRT " << Id << " Direction of [" << targetCoord.x << ":" << targetCoord.y
        //         << "] changed from " << static_cast<int>(src_dir) << " to "
        //         << static_cast<int>(tgt_dir) << log4cpp::eol;
        crt.m_flag = true
	break
    }

    // mutation: constant correction
    for {
        break

        if GetMutator().chanceToConstantChanged() || crt.m_flag {
            break
        }

        opt := value.Options()

        if opt.sourceOpType != Ot_ConstantVector {
            break
        }

        GetMutator().MutateVector(&value.arg2)
        crt.m_flag = true
	break
    }

    // instruction copying
    target.instruction = *value
    target.tailId = crt.Id
    target.executorId = 0
    target.dir = tgt_dir
    crt.embrion = append(crt.embrion, *targetCoord)

    //  mutation: insert instruction
    for {
        if GetMutator().chanceToNopInserted() || crt.m_flag {
            break
        }

        // корректируем указатель на цель, чтобы дальнейшее копирование выполнялось корректно
        cell := GetWorld().GetCellByCoord(crt.ptr)
        vect := &crt.internal_memory[cell.instruction.arg1.x]
        *vect = vect.Next(tgt_dir)

        // вставляем еще одну инструкцию
        to := targetCoord.Next(tgt_dir)
        targetCell := GetWorld().GetCellByCoord(to)
        targetCell.instruction.code = Op_Nop
        targetCell.tailId = crt.Id
        targetCell.executorId = 0
        targetCell.dir = tgt_dir
	crt.embrion = append(crt.embrion, to)

        crt.m_flag = true

        //Log::Not << "CRT " << Id << " <<Nop inserted into [" << to.x << ":" << to.y << "]>>" << log4cpp::eol;
	break

    }

}
