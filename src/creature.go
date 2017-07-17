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
