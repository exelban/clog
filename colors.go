package logg

// Base colors
const (
	Black int = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// High intensity colors
const (
	HiBlack int = iota + 60
	HiRed
	HiGreen
	HiYellow
	HiBlue
	HiMagenta
	HiCyan
	HiWhite
)

// Base colors
type Colors interface {
	Black() string
	Red() string
	Green() string
	Yellow() string
	Blue() string
	Magenta() string
	Cyan() string
	White() string

	HiBlack() string
	HiRed() string
	HiGreen() string
	HiYellow() string
	HiBlue() string
	HiMagenta() string
	HiCyan() string
	HiWhite() string
}

type colors struct{}

// Black text color.
func (c *colors) Black() string { return generate(Black) }

// Red text color.
func (c *colors) Red() string { return generate(Red) }

// Green text color.
func (c *colors) Green() string { return generate(Green) }

// Yellow text color.
func (c *colors) Yellow() string { return generate(Yellow) }

// Blue text color.
func (c *colors) Blue() string { return generate(Blue) }

// Magenta text color.
func (c *colors) Magenta() string { return generate(Magenta) }

// Cyan text color.
func (c *colors) Cyan() string { return generate(Cyan) }

// White text color.
func (c *colors) White() string { return generate(White) }

// Black high intense color.
func (c *colors) HiBlack() string { return generate(HiBlack) }

// Red high intense color.
func (c *colors) HiRed() string { return generate(HiRed) }

// Green high intense color.
func (c *colors) HiGreen() string { return generate(HiGreen) }

// Yellow high intense color.
func (c *colors) HiYellow() string { return generate(HiYellow) }

// Blue high intense color.
func (c *colors) HiBlue() string { return generate(HiBlue) }

// Magenta high intense color.
func (c *colors) HiMagenta() string { return generate(HiMagenta) }

// Cyan high intense color.
func (c *colors) HiCyan() string { return generate(HiCyan) }

// White high intense color.
func (c *colors) HiWhite() string { return generate(HiWhite) }