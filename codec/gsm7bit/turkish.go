package gsm7bit

// 土耳其 gsm7bit 编码
// UDH contains 0x25 0x01 0x01
var GsmTurkish = []rune{
	'@', 'Δ', ' ', '0', 'İ', 'P', 'ç', 'p',
	'£', '_', '!', '1', 'A', 'Q', 'a', 'q',
	'$', 'Φ', '"', '2', 'B', 'R', 'b', 'r',
	'¥', 'Γ', '#', '3', 'C', 'S', 'c', 's',
	'€', 'Λ', '¤', '4', 'D', 'T', 'd', 't',
	'é', 'Ω', '%', '5', 'E', 'U', 'e', 'u',
	'ù', 'Π', '&', '6', 'F', 'V', 'f', 'v',
	'ı', 'Ψ', '\'', '7', 'G', 'W', 'g', 'w',
	'ò', 'Σ', '(', '8', 'H', 'X', 'h', 'x',
	'Ç', 'Θ', ')', '9', 'I', 'Y', 'i', 'y',
	'\n', 'Ξ', '*', ':', 'J', 'Z', 'j', 'z',
	'Ğ', esc, '+', ';', 'K', 'Ä', 'k', 'ä',
	'ğ', 'Ş', ',', '<', 'L', 'Ö', 'l', 'ö',
	'\r', 'ş', '-', '=', 'M', 'Ñ', 'm', 'ñ',
	'Å', 'ß', '.', '>', 'N', 'Ü', 'n', 'ü',
	'å', 'É', '/', '?', 'O', '§', 'o', 'à',
}

// UDH contains 0x24 0x01 0x01
var GsmTurkishExt = []rune{
	zr0, zr0, zr0, zr0, '|', zr0, zr0, zr0,
	zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, zr0, zr0, 'Ş', 'ç', 'ş',
	zr0, '^', zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, zr0, zr0, zr0, '€', zr0,
	zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, zr0, 'Ğ', zr0, 'ğ', zr0,
	zr0, zr0, '{', zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, '}', zr0, 'İ', zr0, 'ı', zr0,
	'\f', zr0, zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, '[', zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, '~', zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, ']', zr0, zr0, zr0, zr0,
	zr0, zr0, '\\', zr0, zr0, zr0, zr0, zr0,
}
