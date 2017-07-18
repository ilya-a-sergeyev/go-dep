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
        cell := TheWorld.GetCellByCoord(crt.fingerprint[i])
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
            cell := TheWorld.GetCellByCoord(crt.embrion[i])
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
    ins := TheWorld.GetCellByCoord(crt.ptr).instruction
    dir := TheWorld.GetCellByCoord(crt.ptr).dir

    for !(ins.code == Op_End) && !(ins.code == Op_None) {
        //Log::Inf << "CRT " << Id << " [" << ptr.x << ":" << ptr.y << "] " << static_cast<int>(dir);
        //Log::Inf << log4cpp::eol;
        crt.ptr = crt.ptr.Next(dir)
        ins = TheWorld.GetCellByCoord(crt.ptr).instruction
        dir = TheWorld.GetCellByCoord(crt.ptr).dir
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
            tcell := TheWorld.GetCellByCoord(crt.ptr)
            crt.ptr.Dec(tcell.dir)
            cnt--
        }
        return
    }
    if steps>0 {
        cnt := steps;
        for cnt>0 {
            tcell := TheWorld.GetCellByCoord(crt.ptr)
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

