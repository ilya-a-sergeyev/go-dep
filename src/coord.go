package godep

type Direction int

const (
	Forward = Direction(iota)
	Back
	Left
	Right
)

type Coord struct {
	x	int
	y	int
}

func (p *Coord) Next(dir Direction) Coord {
	newcoord := Coord { x:p.x, y:p.y }
	switch dir {
		case Forward:
			newcoord.x = (newcoord.x+1) % TheWorld.Size
			// next row
			if newcoord.x == 0 {
				newcoord.y = (newcoord.y+1)%TheWorld.Size
			}
		case Right:
			newcoord.y = (newcoord.y+1) % TheWorld.Size
			// next column
			if newcoord.y == 0 {
				newcoord.x = (newcoord.x+1)%TheWorld.Size
			}
		case Back:
			if (p.x > 0) {
				newcoord.x = p.x-1
			} else {
				newcoord.x = TheWorld.Size-1
			}
		case Left:
			if (p.y > 0) {
				newcoord.y = p.y-1
			} else {
				newcoord.y = TheWorld.Size-1
			}
	}
	return newcoord
}

func (p *Coord) Prev(dir Direction) Coord {
	newcoord := Coord { x:p.x, y:p.y }
	switch dir {
		case Back:
			newcoord.x = (newcoord.x+1) % TheWorld.Size
			// next row
			if newcoord.x == 0 {
				newcoord.y = (newcoord.y+1)%TheWorld.Size
			}
		case Left:
			newcoord.y = (newcoord.y+1) % TheWorld.Size
			// next column
			if newcoord.y == 0 {
				newcoord.x = (newcoord.x+1)%TheWorld.Size
			}
		case Forward:
			if p.x > 0 {
				newcoord.x = p.x-1
			} else {
				newcoord.x = TheWorld.Size-1
			}
		case Right:
			if p.y > 0 {
				newcoord.y = p.y-1
			} else {
				newcoord.y = TheWorld.Size-1
			}
	}
	return newcoord
}

func (p *Coord) Inc(dir Direction) {
	switch dir {
	case Forward: p.x++
	case Back: p.x--
	case Left: p.y--
	case Right: p.y++
	}
}

func (p *Coord) Dec(dir Direction) {
	switch dir {
	case Forward: p.x--
	case Back: p.x++
	case Left: p.y++
	case Right: p.y--
	}
}
