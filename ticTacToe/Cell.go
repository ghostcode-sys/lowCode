package tictactoe

type Cell struct {
	row    int
	col    int
	symbol Symbol
}

func NewCell(r, c int) *Cell {
	return &Cell{
		row:    r,
		col:    c,
		symbol: EMPTY,
	}
}

func (c *Cell) IsEmpty() bool {
	return c.symbol == EMPTY
}

func (c *Cell) SetSymbol(s Symbol) {
	c.symbol = s
}

func (c *Cell) GetSymbol()Symbol {
	return c.symbol
}

func (c *Cell) GetCoordinate() (int, int) {
	return c.row, c.col
}
