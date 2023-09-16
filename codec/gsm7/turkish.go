package gsm7

// 土耳其 gsm7bit 编码
// UDH contains 0x25 0x01 0x01
var (
	gsmTurkish = []rune{
		'@', '£', '$', '¥', '€', 'é', 'ù', 'ı', 'ò', 'Ç', '\n', 'Ğ', 'ğ', '\r', 'Å', 'å',
		'Δ', '_', 'Φ', 'Γ', 'Λ', 'Ω', 'Π', 'Ψ', 'Σ', 'Θ', 'Ξ', esc, 'Ş', 'ş', 'ß', 'É',
		' ', '!', '"', '#', '¤', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':', ';', '<', '=', '>', '?',
		'İ', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
		'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'Ä', 'Ö', 'Ñ', 'Ü', '§',
		'ç', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
		'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'ä', 'ö', 'ñ', 'ü', 'à',
	}

	// UDH contains 0x24 0x01 0x01
	gsmTurkishExt = []rune{
		zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, '\f', zr0, zr0, zr0, zr0, zr0,
		zr0, zr0, zr0, zr0, '^', zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
		zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, '{', '}', zr0, zr0, zr0, zr0, zr0, '\\',
		zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, '[', '~', ']', zr0,
		'|', zr0, zr0, zr0, zr0, zr0, zr0, 'Ğ', zr0, 'İ', zr0, zr0, zr0, zr0, zr0, zr0,
		zr0, zr0, zr0, 'Ş', zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
		zr0, zr0, zr0, 'ç', zr0, '€', zr0, 'ğ', zr0, 'ı', zr0, zr0, zr0, zr0, zr0, zr0,
		zr0, zr0, zr0, 'ş', zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
	}
	gsmDeTurkish    map[byte]rune
	gsmDeTurkishExt map[byte]rune

	gsmEnTurkish    map[rune]byte
	gsmEnTurkishExt map[rune]byte
)

func init() {
	gsmDeTurkish = make(map[byte]rune)
	gsmDeTurkishExt = make(map[byte]rune)
	gsmEnTurkish = make(map[rune]byte)
	gsmEnTurkishExt = make(map[rune]byte)
	for i, v := range gsmTurkish {
		gsmDeTurkish[byte(i)] = v
		gsmEnTurkish[v] = byte(i)
	}
	for i, v := range gsmTurkishExt {
		if v != zr0 {
			gsmDeTurkishExt[byte(i)] = v
			gsmEnTurkishExt[v] = byte(i)
		}
	}
}
