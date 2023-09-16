package gsm7

// 葡萄牙 gsm7bit 编码
// UDH contains 0x25 0x01 0x03
var (
	gsmPortuguese = []rune{
		'@', '£', '$', '¥', 'ê', 'é', 'ú', 'í', 'ó', 'ç', '\n', 'Ô', 'ô', '\r', 'Á', 'á',
		'Δ', '_', 'ª', 'Ç', 'À', '∞', '^', '\\', '€', 'Ó', '|', esc, 'Â', 'â', 'Ê', 'É',
		' ', '!', '"', '#', 'º', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':', ';', '<', '=', '>', '?',
		'Í', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
		'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'Ã', 'Õ', 'Ú', 'Ü', '§',
		'~', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
		'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'ã', 'õ', '`', 'ü', 'à',
	}

	// UDH contains 0x24 0x01 0x03
	gsmPortugueseExt = []rune{
		zr0, zr0, zr0, zr0, zr0, 'ê', zr0, zr0, zr0, 'ç', '\f', 'Ô', 'ô', zr0, 'Á', 'á',
		zr0, zr0, 'Φ', 'Γ', '^', 'Ω', 'Π', 'Ψ', 'Σ', 'Θ', zr0, zr0, zr0, zr0, zr0, 'Ê',
		zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, '{', '}', zr0, zr0, zr0, zr0, zr0, '\\',
		zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, '[', '~', ']', zr0,
		'|', 'À', zr0, zr0, zr0, zr0, zr0, zr0, zr0, 'Í', zr0, zr0, zr0, zr0, zr0, 'Ó',
		zr0, zr0, zr0, zr0, zr0, 'Ú', zr0, zr0, zr0, zr0, zr0, 'Ã', 'Õ', zr0, zr0, zr0,
		zr0, 'Â', zr0, zr0, zr0, '€', zr0, zr0, zr0, 'í', zr0, zr0, zr0, zr0, zr0, 'ó',
		zr0, zr0, zr0, zr0, zr0, 'ú', zr0, zr0, zr0, zr0, zr0, 'ã', 'õ', zr0, zr0, 'â',
	}
	gsmDePortuguese    map[byte]rune
	gsmDePortugueseExt map[byte]rune

	gsmEnPortuguese    map[rune]byte
	gsmEnPortugueseExt map[rune]byte
)

func init() {
	gsmDePortuguese = make(map[byte]rune)
	gsmDePortugueseExt = make(map[byte]rune)
	gsmEnPortuguese = make(map[rune]byte)
	gsmEnPortugueseExt = make(map[rune]byte)
	for i, v := range gsmPortuguese {
		gsmDePortuguese[byte(i)] = v
		gsmEnPortuguese[v] = byte(i)
	}
	for i, v := range gsmPortugueseExt {
		if v != zr0 {
			gsmDePortugueseExt[byte(i)] = v
			gsmEnPortugueseExt[v] = byte(i)
		}
	}
}
