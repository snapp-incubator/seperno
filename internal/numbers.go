package internal

const (
	// Persian Numbers
	faD0 = '\u06F0'
	faD1 = '\u06F1'
	faD2 = '\u06F2'
	faD3 = '\u06F3'
	faD4 = '\u06F4'
	faD5 = '\u06F5'
	faD6 = '\u06F6'
	faD7 = '\u06F7'
	faD8 = '\u06F8'
	faD9 = '\u06F9'

	// Arabic Numbers
	arD0 = '\u0660'
	arD1 = '\u0661'
	arD2 = '\u0662'
	arD3 = '\u0663'
	arD4 = '\u0664'
	arD5 = '\u0665'
	arD6 = '\u0666'
	arD7 = '\u0667'
	arD8 = '\u0668'
	arD9 = '\u0669'

	// English Numbers
	enD0 = '\u0030'
	enD1 = '\u0031'
	enD2 = '\u0032'
	enD3 = '\u0033'
	enD4 = '\u0034'
	enD5 = '\u0035'
	enD6 = '\u0036'
	enD7 = '\u0037'
	enD8 = '\u0038'
	enD9 = '\u0039'

	// Languages for Numbers
	fa = "fa"
	en = "en"
	ar = "ar"
)

var numberConvertor = map[string]map[int32]int32{
	fa: {
		enD0: faD0,
		enD1: faD1,
		enD2: faD2,
		enD3: faD3,
		enD4: faD4,
		enD5: faD5,
		enD6: faD6,
		enD7: faD7,
		enD8: faD8,
		enD9: faD9,
	},
	ar: {
		enD0: arD0,
		enD1: arD1,
		enD2: arD2,
		enD3: arD3,
		enD4: arD4,
		enD5: arD5,
		enD6: arD6,
		enD7: arD7,
		enD8: arD8,
		enD9: arD9,
	},
}

func convertToDestNumber(input int32, lang string) rune {
	switch lang {
	case fa:
		return numberConvertor["fa"][input]
	case ar:
		return numberConvertor["ar"][input]
	default:
		return input
	}
}
