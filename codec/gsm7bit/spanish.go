package gsm7bit

// 西班牙 gsm7bit 编码
var GsmSpanish = []rune{
	'@', 'Δ', ' ', '0', '¡', 'P', '¿', 'p',
	'£', '_', '!', '1', 'A', 'Q', 'a', 'q',
	'$', 'Φ', '"', '2', 'B', 'R', 'b', 'r',
	'¥', 'Γ', '#', '3', 'C', 'S', 'c', 's',
	'è', 'Λ', '¤', '4', 'D', 'T', 'd', 't',
	'é', 'Ω', '%', '5', 'E', 'U', 'e', 'u',
	'ù', 'Π', '&', '6', 'F', 'V', 'f', 'v',
	'ì', 'Ψ', '\'', '7', 'G', 'W', 'g', 'w',
	'ò', 'Σ', '(', '8', 'H', 'X', 'h', 'x',
	'Ç', 'Θ', ')', '9', 'I', 'Y', 'i', 'y',
	'\n', 'Ξ', '*', ':', 'J', 'Z', 'j', 'z',
	'Ø', esc, '+', ';', 'K', 'Ä', 'k', 'ä',
	'ø', 'Æ', ',', '<', 'L', 'Ö', 'l', 'ö',
	'\r', 'æ', '-', '=', 'M', 'Ñ', 'm', 'ñ',
	'Å', 'ß', '.', '>', 'N', 'Ü', 'n', 'ü',
	'å', 'É', '/', '?', 'O', '§', 'o', 'à',
}

// UDH contains 0x24 0x01 0x02
var GsmSpanishExt = []rune{
	zr0, zr0, zr0, zr0, '|', zr0, zr0, zr0,
	zr0, zr0, zr0, zr0, 'Á', zr0, 'á', zr0,
	zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, '^', zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, zr0, zr0, 'Ú', '€', 'ú',
	zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, '{', zr0, zr0, zr0, zr0, zr0,
	'ç', zr0, '}', zr0, 'Í', zr0, 'í', zr0,
	'\f', zr0, zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, '[', zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, '~', zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, ']', zr0, zr0, zr0, zr0,
	zr0, zr0, '\\', zr0, 'Ó', zr0, 'ó', zr0,
}
