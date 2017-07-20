//
// ==========================================
// Vector math for toroidal world
// ==========================================
//

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
			newcoord.x = (newcoord.x+1) % GetWorld().size
			// next row
			if newcoord.x == 0 {
				newcoord.y = (newcoord.y+1)%GetWorld().size
			}
		case Right:
			newcoord.y = (newcoord.y+1) % GetWorld().size
			// next column
			if newcoord.y == 0 {
				newcoord.x = (newcoord.x+1)%GetWorld().size
			}
		case Back:
			if (p.x > 0) {
				newcoord.x = p.x-1
			} else {
				newcoord.x = GetWorld().size-1
			}
		case Left:
			if (p.y > 0) {
				newcoord.y = p.y-1
			} else {
				newcoord.y = GetWorld().size-1
			}
	}
	return newcoord
}

func (p *Coord) Prev(dir Direction) Coord {
	newcoord := Coord { x:p.x, y:p.y }
	switch dir {
		case Back:
			newcoord.x = (newcoord.x+1) % GetWorld().size
			// next row
			if newcoord.x == 0 {
				newcoord.y = (newcoord.y+1)%GetWorld().size
			}
		case Left:
			newcoord.y = (newcoord.y+1) % GetWorld().size
			// next column
			if newcoord.y == 0 {
				newcoord.x = (newcoord.x+1)%GetWorld().size
			}
		case Forward:
			if p.x > 0 {
				newcoord.x = p.x-1
			} else {
				newcoord.x = GetWorld().size-1
			}
		case Right:
			if p.y > 0 {
				newcoord.y = p.y-1
			} else {
				newcoord.y = GetWorld().size-1
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

func (p *Coord) Add(vx, vy int) Coord {

	newcoord := Coord { x:(p.x+vx)%GetWorld().size, y:p.y }

	if p.x+vx >= GetWorld().size {
		newcoord.y = (p.y+vx/GetWorld().size+1)%GetWorld().size
	} else if p.x+vx<0 {
		newcoord.y = (p.y-vx/GetWorld().size-1)%GetWorld().size
		newcoord.x += GetWorld().size
	}

	ty := newcoord.y;
	newcoord.y = (ty+vy)%GetWorld().size

	if p.y+vy >= GetWorld().size {
		newcoord.x = (newcoord.x+vy/GetWorld().size+1)%GetWorld().size
	} else if p.y+vy<0 {
		newcoord.x = (newcoord.x-vy/GetWorld().size-1)%GetWorld().size
		newcoord.y += GetWorld().size
	}

	return newcoord;
}

func (p *Coord) AddRaw(vx, vy int) Coord {
	newcoord := Coord { x:(p.x+vx)%GetWorld().size, y:(p.y+vy)%GetWorld().size }
	return newcoord
}

func sqrtNewton(l int) int {
	var temp, div int

	if l<=0 {
		return 0
	} 

	rslt := l

	if l&0xFFFF0000 != 0 {
		if l&0xFF000000 != 0 {
			div = 0x3FFF
		} else {
			div = 0x3FF
		}
	} else {
		if  l&0x0FF00 != 0 {
			div = 0x3F
		} else {
			if l>4 {
				div = 0x7
			} else {
				div = l
			}
		}
	}

	for true {
		temp = l / div + div
		div =  temp >>  1
		div += temp & 1
		if rslt>div {
			rslt = div
		} else {
			if l/rslt == rslt-1 && l%rslt == 0 {
			rslt--
			}
			return rslt
		}
	}

	return rslt
}

func (p *Coord) Length() int {
	if p.y == 0 {
		if p.x>0 {
			return p.x
		}
		return -p.x
	}
	if p.x == 0 {
		if p.y>0 {
			return p.y
		}
		return -p.y
	}
	return sqrtNewton(p.x*p.x+p.y*p.y)
}
