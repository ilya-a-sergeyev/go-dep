package godep

type Fingerprint []int

type Creature struct {
	Id		int
	entry		Coord
	ptr		Coord
	flags		int

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
    crt.flags = 0
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

    for step := range crt.fingerprint {
        cell := TheWorld.GetCell(step)
        //Log::Inf << "<<!!! CRT " << Id << " " << step%World::worldSize << ":"<< step/World::worldSize << ">>" << log4cpp::eol;
        cell.Clear()
    }

    crt.fingerprint = crt.fingerprint[:0]
    //crt.DoAbort()
}

func (crt *Creature)DoAbort() {
    if len(crt.embrion)>0 {
        //Log::Inf << " CRT " << Id << " has " << embrion.size() << " embrional cells." << log4cpp::eol;
        for step := range crt.embrion {
            cell := TheWorld.GetCell(step)
            if cell.tailId == crt.Id {
                    //Log::Inf <<
                    //    "<<!!! CRT " << Id << " " << step%World::worldSize << ":"<< step/World::worldSize << ">>" << log4cpp::eol;
                cell.Clear()
            }
        }
        crt.embrion = crt.embrion[:0]
    }
}
