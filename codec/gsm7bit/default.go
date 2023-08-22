package gsm7bit

var zr0 = rune(0)
var esc = rune(0x1B)
var data = []rune{ // 7-bit GSM charset
	'@', '£', '$', '¥', 'è', 'é', 'ù', 'ì',
	'ò', 'Ç', '\n', 'Ø', 'ø', '\r', 'Å', 'å',
	'Δ', '_', 'Φ', 'Γ', 'Λ', 'Ω', 'Π', 'Ψ',
	'Σ', 'Θ', 'Ξ', esc, 'Æ', 'æ', 'ß', 'É',
	' ', '!', '"', '#', '¤', '%', '&', '\'',
	'(', ')', '*', '+', ',', '-', '.', '/',
	'0', '1', '2', '3', '4', '5', '6', '7',
	'8', '9', ':', ';', '<', '=', '>', '?',
	'¡', 'A', 'B', 'C', 'D', 'E', 'F', 'G',
	'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
	'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W',
	'X', 'Y', 'Z', 'Ä', 'Ö', 'Ñ', 'Ü', '§',
	'¿', 'a', 'b', 'c', 'd', 'e', 'f', 'g',
	'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
	'p', 'q', 'r', 's', 't', 'u', 'v', 'w',
	'x', 'y', 'z', 'ä', 'ö', 'ñ', 'ü', 'à',
}
var dataExt = []rune{ // 7-bit GSM charset
	zr0, zr0, zr0, zr0, '|', zr0, zr0, zr0,
	zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, '^', zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, zr0, zr0, zr0, '€', zr0,
	zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, '{', zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, '}', zr0, zr0, zr0, zr0, zr0,
	'\f', zr0, zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, '[', zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, '~', zr0, zr0, zr0, zr0,
	zr0, zr0, zr0, ']', zr0, zr0, zr0, zr0,
	zr0, zr0, '\\', zr0, zr0, zr0, zr0, zr0,
}

//'\f': 0x0A, '^': 0x14, '{': 0x28, '}': 0x29, '\\': 0x2F, '[': 0x3C, '~': 0x3D, ']': 0x3E, '|': 0x40, '€': 0x65,
/*
Portuguese language (Latin script)
See also: Portuguese language and Portuguese alphabet
Locking Shift Character Set
for Portuguese language

FF is a Page Break control. If not recognized, it shall be treated like LF.
CR2 is a control character. No language specific character shall be encoded at this position.
SS2 is a second Single Shift Escape control reserved for future extensions.
Turkish language (Latin script)
See also: Turkish language and Turkish alphabet
Locking Shift Character Set
for Turkish language
UDH contains 0x25 0x01 0x01[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70
0x00	@	Δ	SP	0	İ	P	ç	p
0x01	£	_	!	1	A	Q	a	q
0x02	$	Φ	"	2	B	R	b	r
0x03	¥	Γ	#	3	C	S	c	s
0x04	€	Λ	¤	4	D	T	d	t
0x05	é	Ω	%	5	E	U	e	u
0x06	ù	Π	&	6	F	V	f	v
0x07	ı	Ψ	'	7	G	W	g	w
0x08	ò	Σ	(	8	H	X	h	x
0x09	Ç	Θ	)	9	I	Y	i	y
0x0A	LF	Ξ	*	:	J	Z	j	z
0x0B	Ğ	ESC	+	;	K	Ä	k	ä
0x0C	ğ	Ş	,	<	L	Ö	l	ö
0x0D	CR	ş	-	=	M	Ñ	m	ñ
0x0E	Å	ß	.	>	N	Ü	n	ü
0x0F	å	É	/	?	O	§	o	à
LF is a Line Feed control.
CR is a Carriage Return control, or filler.
ESC is an Escape control.
SP is a Space character.
Single Shift Character Set
for Turkish language
UDH contains 0x24 0x01 0x01[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70
0x00	 	 	 	 	|
0x01
0x02
0x03	 	 	 	 	 	Ş	ç	ş
0x04	 	^
0x05	 	 	 	 	 	 	€
0x06
0x07	 	 	 	 	Ğ	 	ğ
0x08	 	 	{
0x09	 	 	}	 	İ	 	ı
0x0A	FF
0x0B	 	SS2
0x0C	 	 	 	[
0x0D	CR2	 	 	~
0x0E	 	 	 	]
0x0F	 	 	\
FF is a Page Break control. If not recognized, it shall be treated like LF.
CR2 is a control character. No language specific character shall be encoded at this position.
SS2 is a second Single Shift Escape control reserved for future extensions.
Urdu language (Arabic and basic Latin scripts)
See also: Urdu language, Urdu alphabet, and Eastern Arabic numerals
It may also be used for the Sindhi language also written in the Arabic script.

Sometimes it may be used for Arabic language as well, but the Eastern digits (encoded here in their Persian-Hindu variant) won't be used in that case because standard Arabic prefer its traditional Eastern Arabic digits, and will frequently be replaced by Western Arabic digits (encoded in the locking shift character set in column 0x30) which are also used now frequently in Urdu as well. However, in India, phones recognizing the Arabic language indication may substitute the Persian-Hindu variants of the Eastern Arabic digits by the traditional Eastern Arabic digits.

Locking Shift Character Set
for Urdu language
UDH contains 0x25 0x01 0x0D[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70
0x00	ا	ث	SP	0	ص	ں	◌ٔ	p
0x01	آ	ج	!	1	ض	ڻ	a	q
0x02	ب	ځ	ڏ	2	ط	ڼ	b	r
0x03	ٻ	ڄ	ڍ	3	ظ	و	c	s
0x04	ڀ	ڃ	ذ	4	ع	ۄ	d	t
0x05	پ	څ	ر	5	ف	ە	e	u
0x06	ڦ	چ	ڑ	6	ق	ہ	f	v
0x07	ت	ڇ	ړ	7	ک	ھ	g	w
0x08	ۂ	ح	)	8	ڪ	ء	h	x
0x09	ٿ	خ	(	9	ګ	ی	i	y
0x0A	LF	د	ڙ	:	گ	ې	j	z
0x0B	ٹ	ESC	ز	;	ڳ	ے	k	◌ٕ
0x0C	ٽ	ڌ	,	ښ	ڱ	◌ٍ	l	◌ّ
0x0D	CR	ڈ	ږ	س	ل	◌ِ	m	◌ٓ
0x0E	ٺ	ډ	.	ش	م	◌ُ	n	◌ٖ
0x0F	ټ	ڊ	ژ	?	ن	◌ٗ	o	◌ٰ
LF is a Line Feed control.
CR is a Carriage Return control, or filler.
ESC is an Escape control.
SP is a Space character.
Single Shift Character Set
for Urdu language
UDH contains 0x24 0x01 0x0D[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70
0x00	@	<	۴	◌ؓ	|	P
0x01	£	=	۵	◌ؔ	A	Q
0x02	$	>	۶	؛	B	R
0x03	¥	¡	۷	؟	C	S
0x04	¿	^	۸	ـ	D	T
0x05	"	¡	۹	◌ْ	E	U	€
0x06	¤	_	،	◌٘	F	V
0x07	%	#	؍	٫	G	W
0x08	&	*	{	٬	H	X
0x09	'	؀	}	ٲ	I	Y
0x0A	FF	؁	؎	ٳ	J	Z
0x0B	*	SS2	؏	ۍ	K
0x0C	+	۰	◌ؐ	[	L
0x0D	CR2	۱	◌ؑ	~	M
0x0E	-	۲	◌ؒ	]	N
0x0F	/	۳	\	۔	O
FF is a Page Break control. If not recognized, it shall be treated like LF.
CR2 is a control character. No language specific character shall be encoded at this position.
SS2 is a second Single Shift Escape control reserved for future extensions.
Hindi language (Devanagari and basic Latin scripts)
See also: Standard Hindi and Devanagari
Locking Shift Character Set
for Hindi language
UDH contains 0x25 0x01 0x06[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70
0x00	◌ँ	ऐ	SP	0	ब	◌ा	ॐ	p
0x01	◌ं	ऑ	!	1	भ	◌ि	a	q
0x02	◌ः	ऒ	ट	2	म	◌ी	b	r
0x03	अ	ओ	ठ	3	य	◌ु	c	s
0x04	आ	औ	ड	4	र	◌ू	d	t
0x05	इ	क	ढ	5	ऱ	◌ृ	e	u
0x06	ई	ख	ण	6	ल	◌ॄ	f	v
0x07	उ	ग	त	7	ळ	◌ॅ	g	w
0x08	ऊ	घ	)	8	ऴ	◌ॆ	h	x
0x09	ऋ	ङ	(	9	व	◌े	i	y
0x0A	LF	च	थ	:	श	◌ै	j	z
0x0B	ऌ	ESC	द	;	ष	◌ॉ	k	ॲ
0x0C	ऍ	छ	,	ऩ	स	◌ॊ	l	ॻ
0x0D	CR	ज	ध	प	ह	◌ो	m	ॼ
0x0E	ऎ	झ	.	फ	◌़	◌ौ	n	ॾ
0x0F	ए	ञ	न	?	ऽ	◌्	o	ॿ
LF is a Line Feed control.
CR is a Carriage Return control, or filler.
ESC is an Escape control.
SP is a Space character.
Single Shift Character Set
for Hindi language
UDH contains 0x24 0x01 0x06[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70
0x00	@	<	४	ज़	|	P
0x01	£	=	५	ड़	A	Q
0x02	$	>	६	ढ़	B	R
0x03	¥	¡	७	फ़	C	S
0x04	¿	^	८	य़	D	T
0x05	"	¡	९	ॠ	E	U	€
0x06	¤	_	◌॑	ॡ	F	V
0x07	%	#	◌॒	◌ॢ	G	W
0x08	&	*	{	◌ॣ	H	X
0x09	'	।	}	॰	I	Y
0x0A	FF	॥	◌॓	ॱ	J	Z
0x0B	*	SS2	◌॔	 	K
0x0C	+	०	क़	[	L
0x0D	CR2	१	ख़	~	M
0x0E	-	२	ग़	]	N
0x0F	/	३	\	 	O
FF is a Page Break control. If not recognized, it shall be treated like LF.
CR2 is a control character. No language specific character shall be encoded at this position.
SS2 is a second Single Shift Escape control reserved for future extensions.
Bengali and Assamese languages (Bengali and basic Latin scripts)
See also: Bengali language, Assamese language, and Bengali alphabet
Locking Shift Character Set
for Bengali and Assamese languages
UDH contains 0x25 0x01 0x04[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70
0x00	◌ঁ	ঐ	SP	0	◌ব	া	ৎ	p
0x01	◌ং	 	!	1	ভ	◌ি	a	q
0x02	◌ঃ	 	ট	2	ম	◌ী	b	r
0x03	অ	ও	ঠ	3	য	◌ু	c	s
0x04	আ	ঔ	ড	4	র	◌ূ	d	t
0x05	ই	ক	ঢ	5	 	◌ৃ	e	u
0x06	ঈ	খ	ণ	6	ল	◌ৄ	f	v
0x07	উ	গ	ত	7	 	 	g	w
0x08	ঊ	ঘ	)	8	 	 	h	x
0x09	ঋ	ঙ	(	9	 	◌ে	i	y
0x0A	LF	চ	থ	:	শ	◌ৈ	j	z
0x0B	ঌ	ESC	দ	;	ষ	 	k	◌ৗ
0x0C	 	ছ	,	 	স	 	l	ড়
0x0D	CR	জ	ধ	প	হ	◌ো	m	ঢ়
0x0E	 	ঝ	.	ফ	◌়	◌ৌ	n	ৰ
0x0F	এ	ঞ	ন	?	ঽ	◌্	o	ৱ
LF is a Line Feed control.
CR is a Carriage Return control, or filler.
ESC is an Escape control.
SP is a Space character.
Single Shift Character Set
for Bengali and Assamese languages
UDH contains 0x24 0x01 0x04[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70
0x00	@	<	৬	৶	|	P
0x01	£	=	৭	৷	A	Q
0x02	$	>	৮	৸	B	R
0x03	¥	¡	৯	৹	C	S
0x04	¿	^	য়	৺	D	T
0x05	"	¡	ৠ	 	E	U	€
0x06	¤	_	ৡ	 	F	V
0x07	%	#	◌ৢ	 	G	W
0x08	&	*	{	 	H	X
0x09	'	০	}	 	I	Y
0x0A	FF	১	◌ৣ	 	J	Z
0x0B	*	SS2	৲	 	K
0x0C	+	২	৳	[	L
0x0D	CR2	৩	৴	~	M
0x0E	-	৪	৵	]	N
0x0F	/	৫	\	 	O
FF is a Page Break control. If not recognized, it shall be treated like LF.
CR2 is a control character. No language specific character shall be encoded at this position.
SS2 is a second Single Shift Escape control reserved for future extensions.
Punjabi language (Gurmukhī and basic Latin scripts)
See also: Punjabi language and Gurmukhī alphabet
Locking Shift Character Set
for Punjabi language
UDH contains 0x25 0x01 0x0A[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70
0x00	◌ਁ	ਐ	SP	0	ਬ	◌ਾ	◌ੑ	p
0x01	◌ਂ	 	!	1	ਭ	◌ਿ	a	q
0x02	◌ਃ	 	ਟ	2	ਮ	◌ੀ	b	r
0x03	ਅ	ਓ	ਠ	3	ਯ	◌ੁ	c	s
0x04	ਆ	ਔ	ਡ	4	ਰ	◌ੂ	d	t
0x05	ਇ	ਕ	ਢ	5	 	 	e	u
0x06	ਈ	ਖ	ਣ	6	ਲ	 	f	v
0x07	ਉ	ਗ	ਤ	7	ਲ਼	 	g	w
0x08	ਊ	ਘ	)	8	 	 	h	x
0x09	 	ਙ	(	9	ਵ	◌ੇ	i	y
0x0A	LF	ਚ	ਥ	:	ਸ਼	◌ੈ	j	z
0x0B	 	ESC	ਦ	;	 	 	k	◌ੰ
0x0C	 	ਛ	,	 	ਸ	 	l	◌ੱ
0x0D	CR	ਜ	ਧ	ਪ	ਹ	◌ੋ	m	ੲ
0x0E	 	ਝ	.	ਫ	◌਼	◌ੌ	n	ੳ
0x0F	ਏ	ਞ	ਨ	?	 	◌੍	o	ੴ
LF is a Line Feed control.
CR is a Carriage Return control, or filler.
ESC is an Escape control.
SP is a Space character.
Single Shift Character Set
for Punjabi language
UDH contains 0x24 0x01 0x0A[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70
0x00	@	<	੪	 	|	P
0x01	£	=	੫	 	A	Q
0x02	$	>	੬	 	B	R
0x03	¥	¡	੭	 	C	S
0x04	¿	^	੮	 	D	T
0x05	"	¡	੯	 	E	U	€
0x06	¤	_	ਖ਼	 	F	V
0x07	%	#	ਗ਼	 	G	W
0x08	&	*	{	 	H	X
0x09	'	।	}	 	I	Y
0x0A	FF	॥	ਜ਼	 	J	Z
0x0B	*	SS2	ੜ	 	K
0x0C	+	੦	ਫ਼	[	L
0x0D	CR2	੧	◌ੵ	~	M
0x0E	-	੨	 	]	N
0x0F	/	੩	\	 	O
FF is a Page Break control. If not recognized, it shall be treated like LF.
CR2 is a control character. No language specific character shall be encoded at this position.
SS2 is a second Single Shift Escape control reserved for future extensions.
Gujarati language (Gujarati and basic Latin scripts)
See also: Gujarati language and Gujarati alphabet
Locking Shift Character Set
for Gujarati language
UDH contains 0x25 0x01 0x05[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70
0x00	◌ઁ	ઐ	SP	0	બ	◌ા	ૐ	p
0x01	◌ં	ઑ	!	1	ભ	◌િ	a	q
0x02	◌ઃ	 	ટ	2	મ	◌ી	b	r
0x03	અ	ઓ	ઠ	3	ય	◌ુ	c	s
0x04	આ	ઔ	ડ	4	ર	◌ૂ	d	t
0x05	ઇ	ક	ઢ	5	 	◌ૃ	e	u
0x06	ઈ	ખ	ણ	6	લ	◌ૄ	f	v
0x07	ઉ	ગ	ત	7	ળ	◌ૅ	g	w
0x08	ઊ	ઘ	)	8	 	 	h	x
0x09	ઋ	ઙ	(	9	વ	◌ે	i	y
0x0A	LF	ચ	થ	:	શ	◌ૈ	j	z
0x0B	ઌ	ESC	દ	;	ષ	◌ૉ	k	ૠ
0x0C	ઍ	છ	,	 	સ	 	l	ૡ
0x0D	CR	જ	ધ	પ	હ	◌ો	m	◌ૢ
0x0E	 	ઝ	.	ફ	◌઼	◌ૌ	n	◌ૣ
0x0F	એ	ઞ	ન	?	ઽ	◌્	o	૱
LF is a Line Feed control.
CR is a Carriage Return control, or filler.
ESC is an Escape control.
SP is a Space character.
Single Shift Character Set
for Gujarati language
UDH contains 0x24 0x01 0x05[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70
0x00	@	<	૪	 	|	P
0x01	£	=	૫	 	A	Q
0x02	$	>	૬	 	B	R
0x03	¥	¡	૭	 	C	S
0x04	¿	^	૮	 	D	T
0x05	"	¡	૯	 	E	U	€
0x06	¤	_	 	 	F	V
0x07	%	#	 	 	G	W
0x08	&	*	{	 	H	X
0x09	'	।	}	 	I	Y
0x0A	FF	॥	 	 	J	Z
0x0B	*	SS2	 	 	K
0x0C	+	૦	 	[	L
0x0D	CR2	૧	 	~	M
0x0E	-	૨	 	]	N
0x0F	/	૩	\	 	O
FF is a Page Break control. If not recognized, it shall be treated like LF.
CR2 is a control character. No language specific character shall be encoded at this position.
SS2 is a second Single Shift Escape control reserved for future extensions.
Oriya language (Oriya and basic Latin scripts)
See also: Oriya language and Oriya alphabet
Locking Shift Character Set
for Oriya language
UDH contains 0x25 0x01 0x09[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70
0x00	◌ଁ	ଐ	SP	0	ବ	◌ା	◌ୖ	p
0x01	◌ଂ	 	!	1	ଭ	◌ି	a	q
0x02	◌ଃ	 	ଟ	2	ମ	◌ୀ	b	r
0x03	ଅ	ଓ	ଠ	3	ଯ	◌ୁ	c	s
0x04	ଆ	ଔ	ଡ	4	ର	◌ୂ	d	t
0x05	ଇ	କ	ଢ	5	 	◌ୃ	e	u
0x06	ଈ	ଖ	ଣ	6	ଲ	ୄ	f	v
0x07	ଉ	ଗ	ତ	7	ଳ	 	g	w
0x08	ଊ	ଘ	)	8	 	 	h	x
0x09	ଋ	ଙ	(	9	ଵ	◌େ	i	y
0x0A	LF	ଚ	ଥ	:	ଶ	◌ୈ	j	z
0x0B	ଌ	ESC	ଦ	;	ଷ	 	k	◌ୗ
0x0C	 	ଛ	,	 	ସ	 	l	ୠ
0x0D	CR	ଜ	ଧ	ପ	ହ	◌ୋ	m	ୡ
0x0E	 	ଝ	.	ଫ	◌଼	◌ୌ	n	◌ୢ
0x0F	ଏ	ଞ	ନ	?	ଽ	◌୍	o	◌ୣ
LF is a Line Feed control.
CR is a Carriage Return control, or filler.
ESC is an Escape control.
SP is a Space character.
Single Shift Character Set
for Oriya language
UDH contains 0x24 0x01 0x09[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70
0x00	@	<	୪	 	|	P
0x01	£	=	୫	 	A	Q
0x02	$	>	୬	 	B	R
0x03	¥	¡	୭	 	C	S
0x04	¿	^	୮	 	D	T
0x05	"	¡	୯	 	E	U	€
0x06	¤	_	ଡ଼	 	F	V
0x07	%	#	ଢ଼	 	G	W
0x08	&	*	{	 	H	X
0x09	'	।	}	 	I	Y
0x0A	FF	॥	ୟ	 	J	Z
0x0B	*	SS2	୰	 	K
0x0C	+	୦	ୱ	[	L
0x0D	CR2	୧	 	~	M
0x0E	-	୨	 	]	N
0x0F	/	୩	\	 	O
FF is a Page Break control. If not recognized, it shall be treated like LF.
CR2 is a control character. No language specific character shall be encoded at this position.
SS2 is a second Single Shift Escape control reserved for future extensions.
Tamil language (Tamil and basic Latin scripts)
See also: Tamil language and Tamil alphabet
Locking Shift Character Set
for Tamil language
UDH contains 0x25 0x01 0x0B[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70
0x00	 	ஐ	SP	0	 	◌ா	ௐ	p
0x01	◌ஂ	 	!	1	 	◌ி	a	q
0x02	◌ஃ	ஒ	ட	2	ம	◌ீ	b	r
0x03	அ	ஓ	 	3	ய	◌ு	c	s
0x04	ஆ	ஔ	 	4	ர	◌ூ	d	t
0x05	இ	க	 	5	ற	 	e	u
0x06	ஈ	 	ண	6	ல	 	f	v
0x07	உ	 	த	7	ள	 	g	w
0x08	ஊ	 	)	8	ழ	◌ெ	h	x
0x09	 	ங	(	9	வ	◌ே	i	y
0x0A	LF	ச	 	:	ஶ	◌ை	j	z
0x0B	 	ESC	 	;	ஷ	 	k	◌ௗ
0x0C	 	 	,	ன	ஸ	◌ொ	l	௰
0x0D	CR	ஜ	 	ப	ஹ	◌ோ	m	௱
0x0E	எ	 	.	 	 	◌ௌ	n	௲
0x0F	ஏ	ஞ	ந	?	 	◌்	o	௹
LF is a Line Feed control.
CR is a Carriage Return control, or filler.
ESC is an Escape control.
SP is a Space character.
Single Shift Character Set
for Tamil language
UDH contains 0x24 0x01 0x0B[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70
0x00	@	<	௪	 	|	P
0x01	£	=	௫	 	A	Q
0x02	$	>	௬	 	B	R
0x03	¥	¡	௭	 	C	S
0x04	¿	^	௮	 	D	T
0x05	"	¡	௯	 	E	U	€
0x06	¤	_	௳	 	F	V
0x07	%	#	௴	 	G	W
0x08	&	*	{	 	H	X
0x09	'	।	}	 	I	Y
0x0A	FF	॥	௵	 	J	Z
0x0B	*	SS2	௶	 	K
0x0C	+	௦	௷	[	L
0x0D	CR2	௧	௸	~	M
0x0E	-	௨	௺	]	N
0x0F	/	௩	\	 	O
FF is a Page Break control. If not recognized, it shall be treated like LF.
CR2 is a control character. No language specific character shall be encoded at this position.
SS2 is a second Single Shift Escape control reserved for future extensions.
Telugu language (Telugu and basic Latin scripts)
See also: Telugu language and Telugu alphabet
Locking Shift Character Set
for Telugu language
UDH contains 0x25 0x01 0x0C[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70
0x00	◌ఁ	ఐ	SP	0	బ	◌ా	◌ౕ	p
0x01	◌ం	 	!	1	భ	◌ి	a	q
0x02	◌ః	ఒ	ట	2	మ	◌ీ	b	r
0x03	అ	ఓ	ఠ	3	య	◌ు	c	s
0x04	ఆ	ఔ	డ	4	ర	◌ూ	d	t
0x05	ఇ	క	ఢ	5	ఱ	◌ృ	e	u
0x06	ఈ	ఖ	ణ	6	ల	◌ౄ	f	v
0x07	ఉ	గ	త	7	ళ	 	g	w
0x08	ఊ	ఘ	)	8	 	◌ె	h	x
0x09	ఋ	ఙ	(	9	వ	◌ే	i	y
0x0A	LF	చ	థ	:	శ	◌ై	j	z
0x0B	ఌ	ESC	ద	;	ష	 	k	◌ౖ
0x0C	 	ఛ	,	 	స	◌ొ	l	ౠ
0x0D	CR	జ	ధ	ప	హ	◌ో	m	ౡ
0x0E	ఎ	ఝ	.	ఫ	 	◌ౌ	n	◌ౢ
0x0F	ఏ	ఞ	న	?	ఽ	◌్	o	◌ౣ
LF is a Line Feed control.
CR is a Carriage Return control, or filler.
ESC is an Escape control.
SP is a Space character.
Single Shift Character Set
for Telugu language
UDH contains 0x24 0x01 0x0C[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70
0x00	@	<	౪	౽	|	P
0x01	£	=	౫	౾	A	Q
0x02	$	>	౬	౿	B	R
0x03	¥	¡	౭	 	C	S
0x04	¿	^	౮	 	D	T
0x05	"	¡	౯	 	E	U
0x06	¤	_	ౘ	 	F	V
0x07	%	#	ౙ	 	G	W
0x08	&	*	{	 	H	X
0x09	'	 	}	 	I	Y
0x0A	FF	 	౸	 	J	Z
0x0B	*	SS2	౹	 	K
0x0C	+	౦	౺	[	L
0x0D	CR2	౧	౻	~	M
0x0E	-	౨	౼	]	N
0x0F	/	౩	\	 	O
FF is a Page Break control. If not recognized, it shall be treated like LF.
CR2 is a control character. No language specific character shall be encoded at this position.
SS2 is a second Single Shift Escape control reserved for future extensions.
Kannada language (Kannada and basic Latin scripts)
See also: Kannada language and Kannada alphabet
Locking Shift Character Set
for Kannada language
UDH contains 0x25 0x01 0x07[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70 !
0x00	 	ಐ	SP	0	ಬ	ಾ	ೕ	p
0x01	ಂ	 	!	1	ಭ	ಿ	a	q
0x02	ಃ	ಒ	ಟ	2	ಮ	ೀ	b	r
0x03	ಅ	ಓ	ಠ	3	ಯ	ು	c	s
0x04	ಆ	ಔ	ಪ	4	ರ	ೂ	d	t
0x05	ಇ	ಕ	ಢ	5	ಱ	ೃ	e	u
0x06	ಈ	ಖ	ಣ	6	ಲ	ೄ	f	v
0x07	ಉ	ಗ	ತ	7	ಳ	 	g	w
0x08	ಊ	ಘ	)	8	 	ೆ	h	x
0x09	ಋ	ಙ	(	9	ವ	ೇ	i	y
0x0A	LF	ಚ	ಥ	:	ಶ	ೈ	j	z
0x0B	ಌ	ESC	ದ	;	ಷ	 	k	ೖ
0x0C	 	ಛ	,	 	ಸ	ೊ	l	ೠ
0x0D	CR	ಜ	ಧ	ಪ	ಹ	ೋ	m	ೡ
0x0E	ಎ	ಝ	.	ಫ	಼	ೌ	n	ೢ
0x0F	ಏ	ಞ	ನ	?	ಽ	್	o	ೣ
LF is a Line Feed control.
CR is a Carriage Return control, or filler.
ESC is an Escape control.
SP is a Space character.
Single Shift Character Set
for Kannada language
UDH contains 0x24 0x01 0x07[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70 !
0x00	@	<	೪	 	|	P
0x01	£	=	೫	 	A	Q
0x02	$	>	೬	 	B	R
0x03	¥	¡	೭	 	C	S
0x04	¿	^	೮	 	D	T
0x05	"	¡	೯	 	E	U	€
0x06	¤	_	ೞ	 	F	V
0x07	%	#	ೱ	 	G	W
0x08	&	*	{	 	H	X
0x09	'	।	}	 	I	Y
0x0A	FF	॥	ೲ	 	J	Z
0x0B	*	SS2	 	 	K
0x0C	+	೦	 	]	L
0x0D	CR2	೧	 	~	M
0x0E	-	೨	 	]	N
0x0F	/	೩	\	 	O
FF is a Page Break control. If not recognized, it shall be treated like LF.
CR2 is a control character. No language specific character shall be encoded at this position.
SS2 is a second Single Shift Escape control reserved for future extensions.
Malayalam language (Malayalam and basic Latin scripts)
See also: Malayalam language and Malayalam alphabet
Locking Shift Character Set
for Malayalam language
UDH contains 0x25 0x01 0x08[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70 !
0x00	 	ഐ	SP	0	ബ	ാ	ൗ	p
0x01	ം	 	!	1	ഭ	ി	a	q
0x02	ഃ	ഒ	ട	2	മ	ീ	b	r
0x03	അ	ഓ	ഠ	3	യ	ു	c	s
0x04	ആ	ഔ	ഡ	4	ര	ൂ	d	t
0x05	ഇ	ക	ഢ	5	റ	ൃ	e	u
0x06	ഈ	ഖ	ണ	6	ല	ൄ	f	v
0x07	ഉ	ഗ	ത	7	ള	 	g	w
0x08	ഊ	ഘ	)	8	ഴ	െ	h	x
0x09	ഋ	ങ	(	9	വ	േ	i	y
0x0A	LF	ച	ഥ	:	ശ	ൈ	j	z
0x0B	ഌ	ESC	ദ	;	ഷ	 	k	ൠ
0x0C	 	ഛ	,	 	സ	ൊ	l	ൡ
0x0D	CR	ജ	ധ	പ	ഹ	ോ	m	ൢ
0x0E	എ	ഝ	.	ഫ	 	ൌ	n	ൣ
0x0F	ഏ	ഞ	ന	?	ഽ	്	o	൹
LF is a Line Feed control.
CR is a Carriage Return control, or filler.
ESC is an Escape control.
SP is a Space character.
Single Shift Character Set
for Malayalam language
UDH contains 0x25 0x01 0x08[2]
 	0x00	0x10	0x20	0x30	0x40	0x50	0x60	0x70 !
0x00	@	<	൪	ൻ	-	P
0x01	£	=	൫	ർ	A	Q
0x02	$	>	൬	ൽ	B	R
0x03	¥	¡	൭	ൾ	C	S
0x04	¿	^	൮	ൿ	D	T
0x05	"	¡	൯	 	E	U	€
0x06	¤	_	൰	 	F	V
0x07	%	#	൱	 	G	W
0x08	&	*	{	 	H	X
0x09	'	।	}	 	I	Y
0x0A	FF	॥	൲	 	J	Z
0x0B	*	SS2	൳	 	K
0x0C	+	൦	൴	[	L
0x0D	CR2	൧	൵	~	M
0x0E	-	൨	ൺ	]	N
0x0F	/	൩	\	 	O
*/
