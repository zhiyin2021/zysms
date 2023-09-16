package gsm7

// 西班牙 gsm7bit 编码
var (
	gsmSpanish = []rune{
		'@', '£', '$', '¥', 'è', 'é', 'ù', 'ì', 'ò', 'Ç', '\n', 'Ø', 'ø', '\r', 'Å', 'å',
		'Δ', '_', 'Φ', 'Γ', 'Λ', 'Ω', 'Π', 'Ψ', 'Σ', 'Θ', 'Ξ', esc, 'Æ', 'æ', 'ß', 'É',
		' ', '!', '"', '#', '¤', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':', ';', '<', '=', '>', '?',
		'¡', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
		'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'Ä', 'Ö', 'Ñ', 'Ü', '§',
		'¿', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
		'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'ä', 'ö', 'ñ', 'ü', 'à',
	}

	// UDH contains 0x24 0x01 0x02
	gsmSpanishExt = []rune{
		zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, 'ç', '\f', zr0, zr0, zr0, zr0, zr0,
		zr0, zr0, zr0, zr0, '^', zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
		zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, '{', '}', zr0, zr0, zr0, zr0, zr0, '\\',
		zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, '[', '~', ']', zr0,
		'|', 'Á', zr0, zr0, zr0, zr0, zr0, zr0, zr0, 'Í', zr0, zr0, zr0, zr0, zr0, 'Ó',
		zr0, zr0, zr0, zr0, zr0, 'Ú', zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
		zr0, 'á', zr0, zr0, zr0, '€', zr0, zr0, zr0, 'í', zr0, zr0, zr0, zr0, zr0, 'ó',
	}
	gsmDeSpanish    map[byte]rune
	gsmDeSpanishExt map[byte]rune

	gsmEnSpanish    map[rune]byte
	gsmEnSpanishExt map[rune]byte
)

func init() {
	gsmDeSpanish = make(map[byte]rune)
	gsmDeSpanishExt = make(map[byte]rune)
	gsmEnSpanish = make(map[rune]byte)
	gsmEnSpanishExt = make(map[rune]byte)
	for i, v := range gsmSpanish {
		gsmDeSpanish[byte(i)] = v
		gsmEnSpanish[v] = byte(i)
	}
	for i, v := range gsmSpanishExt {
		if v != zr0 {
			gsmDeSpanishExt[byte(i)] = v
			gsmEnSpanishExt[v] = byte(i)
		}
	}
}
