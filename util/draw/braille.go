package draw

var (
	translationMap = []uint8{0, 3, 1, 4, 2, 5, 6, 7}
)

type BrailleChar uint8

func (c *BrailleChar) Get(x, y uint) bool {
	if x >= 2 || y >= 4 {
		return false
	}
	i := translationMap[y*2+x]
	return uint8(*c)>>i&uint8(1) == uint8(1)
}

func (c *BrailleChar) Set(x, y uint8, active bool) {
	if x >= 2 || y >= 4 {
		return
	}
	i := translationMap[y*2+x]
	if active {
		*c = BrailleChar(uint8(*c) | uint8(1)<<i)
	} else {
		*c = BrailleChar(uint8(*c) & ^(uint8(1) << i))
	}
}

func (c *BrailleChar) Char() string {
	return string(rune(0x2800 + uint16(*c)))
}

type Pos struct {
	X uint
	Y uint
}

type BrailleImage map[Pos]BrailleChar

func (img *BrailleImage) Set(x, y uint, active bool) {
	key := Pos{X: x / 2, Y: y / 4}
	char := (*img)[key]
	char.Set(uint8(x%2), uint8(y%4), active)
	(*img)[key] = char
}

func (img *BrailleImage) Print() {
	var currentPos Pos
	for pos, char := range *img {
		relX := int(pos.X) - int(currentPos.X)
		relY := int(pos.Y) - int(currentPos.Y)
		MoveRel(int(relX), int(relY))
		Write(char.Char())
		currentPos.X = uint(int(currentPos.X) + relX + 1)
		currentPos.Y = uint(int(currentPos.Y) + relY)
	}
	MoveRel(-int(currentPos.X), -int(currentPos.Y))
}
