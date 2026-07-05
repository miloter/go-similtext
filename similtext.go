// Copyright 2026 miloter. All rights reserved.
// Use of this source code is governed by a MIT license which is
// located in the file LICENSE.md.

// Provee métodos para comprobar la similitud entre textos. Particularmente provee la
// posibilidad de realizar comparaciones entre textos sin diferenciar mayúsculas de
// minúscualas y/o los signos diacríticos. También provee métodos para la normalización
// de textos y la comparación por similitud.
package similtext

import (
	"io"
	"math"
	"slices"
	"strings"
	"unicode/utf8"
)

const (
	// En las similitudes de texto, porcentaje mínimo para no sufrir
	// una penalización
	similMinPercent = 50.0
	// En las similitudes de texto, factor de penalización
	similFactorPen = 0.89
	// En las similitudes de texto, factor de penalización cuando
	// hay un número diferente de palabras en cada cadena.
	similNumWordsFactorPen = 0.99
)

type SimilText struct {
	parent map[int32]rune
}

var alphaCode = map[rune]int32{
	'A': 65,
	'B': 66,
	'C': 67,
	'D': 68,
	'E': 69,
	'F': 70,
	'G': 71,
	'H': 72,
	'I': 73,
	'J': 74,
	'K': 75,
	'L': 76,
	'M': 77,
	'N': 78,
	'O': 79,
	'P': 80,
	'Q': 81,
	'R': 82,
	'S': 83,
	'T': 84,
	'U': 85,
	'V': 86,
	'W': 87,
	'X': 88,
	'Y': 89,
	'Z': 90,
	'a': 65,
	'b': 66,
	'c': 67,
	'd': 68,
	'e': 69,
	'f': 70,
	'g': 71,
	'h': 72,
	'i': 73,
	'j': 74,
	'k': 75,
	'l': 76,
	'm': 77,
	'n': 78,
	'o': 79,
	'p': 80,
	'q': 81,
	'r': 82,
	's': 83,
	't': 84,
	'u': 85,
	'v': 86,
	'w': 87,
	'x': 88,
	'y': 89,
	'z': 90,
	'À': 65,
	'Á': 65,
	'Â': 65,
	'Ã': 65,
	'Ä': 65,
	'Å': 65,
	'Ç': 67,
	'È': 69,
	'É': 69,
	'Ê': 69,
	'Ë': 69,
	'Ì': 73,
	'Í': 73,
	'Î': 73,
	'Ï': 73,
	'Ñ': 78,
	'Ò': 79,
	'Ó': 79,
	'Ô': 79,
	'Õ': 79,
	'Ö': 79,
	'Ù': 85,
	'Ú': 85,
	'Û': 85,
	'Ü': 85,
	'Ý': 89,
	'à': 65,
	'á': 65,
	'â': 65,
	'ã': 65,
	'ä': 65,
	'å': 65,
	'ç': 67,
	'è': 69,
	'é': 69,
	'ê': 69,
	'ë': 69,
	'ì': 73,
	'í': 73,
	'î': 73,
	'ï': 73,
	'ñ': 78,
	'ò': 79,
	'ó': 79,
	'ô': 79,
	'õ': 79,
	'ö': 79,
	'ù': 85,
	'ú': 85,
	'û': 85,
	'ü': 85,
	'ý': 89,
	'ÿ': 89,
	'Ā': 65,
	'ā': 65,
	'Ă': 65,
	'ă': 65,
	'Ą': 65,
	'ą': 65,
	'Ć': 67,
	'ć': 67,
	'Ĉ': 67,
	'ĉ': 67,
	'Ċ': 67,
	'ċ': 67,
	'Č': 67,
	'č': 67,
	'Ď': 68,
	'ď': 68,
	'Ē': 69,
	'ē': 69,
	'Ĕ': 69,
	'ĕ': 69,
	'Ė': 69,
	'ė': 69,
	'Ę': 69,
	'ę': 69,
	'Ě': 69,
	'ě': 69,
	'Ĝ': 71,
	'ĝ': 71,
	'Ğ': 71,
	'ğ': 71,
	'Ġ': 71,
	'ġ': 71,
	'Ģ': 71,
	'ģ': 71,
	'Ĥ': 72,
	'ĥ': 72,
	'Ĩ': 73,
	'ĩ': 73,
	'Ī': 73,
	'ī': 73,
	'Ĭ': 73,
	'ĭ': 73,
	'Į': 73,
	'į': 73,
	'İ': 73,
	'Ĵ': 74,
	'ĵ': 74,
	'Ķ': 75,
	'ķ': 75,
	'Ĺ': 76,
	'ĺ': 76,
	'Ļ': 76,
	'ļ': 76,
	'Ľ': 76,
	'ľ': 76,
	'Ń': 78,
	'ń': 78,
	'Ņ': 78,
	'ņ': 78,
	'Ň': 78,
	'ň': 78,
	'Ō': 79,
	'ō': 79,
	'Ŏ': 79,
	'ŏ': 79,
	'Ő': 79,
	'ő': 79,
	'Ŕ': 82,
	'ŕ': 82,
	'Ŗ': 82,
	'ŗ': 82,
	'Ř': 82,
	'ř': 82,
	'Ś': 83,
	'ś': 83,
	'Ŝ': 83,
	'ŝ': 83,
	'Ş': 83,
	'ş': 83,
	'Š': 83,
	'š': 83,
	'Ţ': 84,
	'ţ': 84,
	'Ť': 84,
	'ť': 84,
	'Ũ': 85,
	'ũ': 85,
	'Ū': 85,
	'ū': 85,
	'Ŭ': 85,
	'ŭ': 85,
	'Ů': 85,
	'ů': 85,
	'Ű': 85,
	'ű': 85,
	'Ų': 85,
	'ų': 85,
	'Ŵ': 87,
	'ŵ': 87,
	'Ŷ': 89,
	'ŷ': 89,
	'Ÿ': 89,
	'Ź': 90,
	'ź': 90,
	'Ż': 90,
	'ż': 90,
	'Ž': 90,
	'ž': 90,
	'Ơ': 79,
	'ơ': 79,
	'Ư': 85,
	'ư': 85,
	'Ǎ': 65,
	'ǎ': 65,
	'Ǐ': 73,
	'ǐ': 73,
	'Ǒ': 79,
	'ǒ': 79,
	'Ǔ': 85,
	'ǔ': 85,
	'Ǖ': 85,
	'ǖ': 85,
	'Ǘ': 85,
	'ǘ': 85,
	'Ǚ': 85,
	'ǚ': 85,
	'Ǜ': 85,
	'ǜ': 85,
	'Ǟ': 65,
	'ǟ': 65,
	'Ǡ': 65,
	'ǡ': 65,
	'Ǧ': 71,
	'ǧ': 71,
	'Ǩ': 75,
	'ǩ': 75,
	'Ǫ': 79,
	'ǫ': 79,
	'Ǭ': 79,
	'ǭ': 79,
	'ǰ': 74,
	'Ǵ': 71,
	'ǵ': 71,
	'Ǹ': 78,
	'ǹ': 78,
	'Ǻ': 65,
	'ǻ': 65,
	'Ȁ': 65,
	'ȁ': 65,
	'Ȃ': 65,
	'ȃ': 65,
	'Ȅ': 69,
	'ȅ': 69,
	'Ȇ': 69,
	'ȇ': 69,
	'Ȉ': 73,
	'ȉ': 73,
	'Ȋ': 73,
	'ȋ': 73,
	'Ȍ': 79,
	'ȍ': 79,
	'Ȏ': 79,
	'ȏ': 79,
	'Ȑ': 82,
	'ȑ': 82,
	'Ȓ': 82,
	'ȓ': 82,
	'Ȕ': 85,
	'ȕ': 85,
	'Ȗ': 85,
	'ȗ': 85,
	'Ș': 83,
	'ș': 83,
	'Ț': 84,
	'ț': 84,
	'Ȟ': 72,
	'ȟ': 72,
	'Ȧ': 65,
	'ȧ': 65,
	'Ȩ': 69,
	'ȩ': 69,
	'Ȫ': 79,
	'ȫ': 79,
	'Ȭ': 79,
	'ȭ': 79,
	'Ȯ': 79,
	'ȯ': 79,
	'Ȱ': 79,
	'ȱ': 79,
	'Ȳ': 89,
	'ȳ': 89,
	'Ḁ': 65,
	'ḁ': 65,
	'Ḃ': 66,
	'ḃ': 66,
	'Ḅ': 66,
	'ḅ': 66,
	'Ḇ': 66,
	'ḇ': 66,
	'Ḉ': 67,
	'ḉ': 67,
	'Ḋ': 68,
	'ḋ': 68,
	'Ḍ': 68,
	'ḍ': 68,
	'Ḏ': 68,
	'ḏ': 68,
	'Ḑ': 68,
	'ḑ': 68,
	'Ḓ': 68,
	'ḓ': 68,
	'Ḕ': 69,
	'ḕ': 69,
	'Ḗ': 69,
	'ḗ': 69,
	'Ḙ': 69,
	'ḙ': 69,
	'Ḛ': 69,
	'ḛ': 69,
	'Ḝ': 69,
	'ḝ': 69,
	'Ḟ': 70,
	'ḟ': 70,
	'Ḡ': 71,
	'ḡ': 71,
	'Ḣ': 72,
	'ḣ': 72,
	'Ḥ': 72,
	'ḥ': 72,
	'Ḧ': 72,
	'ḧ': 72,
	'Ḩ': 72,
	'ḩ': 72,
	'Ḫ': 72,
	'ḫ': 72,
	'Ḭ': 73,
	'ḭ': 73,
	'Ḯ': 73,
	'ḯ': 73,
	'Ḱ': 75,
	'ḱ': 75,
	'Ḳ': 75,
	'ḳ': 75,
	'Ḵ': 75,
	'ḵ': 75,
	'Ḷ': 76,
	'ḷ': 76,
	'Ḹ': 76,
	'ḹ': 76,
	'Ḻ': 76,
	'ḻ': 76,
	'Ḽ': 76,
	'ḽ': 76,
	'Ḿ': 77,
	'ḿ': 77,
	'Ṁ': 77,
	'ṁ': 77,
	'Ṃ': 77,
	'ṃ': 77,
	'Ṅ': 78,
	'ṅ': 78,
	'Ṇ': 78,
	'ṇ': 78,
	'Ṉ': 78,
	'ṉ': 78,
	'Ṋ': 78,
	'ṋ': 78,
	'Ṍ': 79,
	'ṍ': 79,
	'Ṏ': 79,
	'ṏ': 79,
	'Ṑ': 79,
	'ṑ': 79,
	'Ṓ': 79,
	'ṓ': 79,
	'Ṕ': 80,
	'ṕ': 80,
	'Ṗ': 80,
	'ṗ': 80,
	'Ṙ': 82,
	'ṙ': 82,
	'Ṛ': 82,
	'ṛ': 82,
	'Ṝ': 82,
	'ṝ': 82,
	'Ṟ': 82,
	'ṟ': 82,
	'Ṡ': 83,
	'ṡ': 83,
	'Ṣ': 83,
	'ṣ': 83,
	'Ṥ': 83,
	'ṥ': 83,
	'Ṧ': 83,
	'ṧ': 83,
	'Ṩ': 83,
	'ṩ': 83,
	'Ṫ': 84,
	'ṫ': 84,
	'Ṭ': 84,
	'ṭ': 84,
	'Ṯ': 84,
	'ṯ': 84,
	'Ṱ': 84,
	'ṱ': 84,
	'Ṳ': 85,
	'ṳ': 85,
	'Ṵ': 85,
	'ṵ': 85,
	'Ṷ': 85,
	'ṷ': 85,
	'Ṹ': 85,
	'ṹ': 85,
	'Ṻ': 85,
	'ṻ': 85,
	'Ṽ': 86,
	'ṽ': 86,
	'Ṿ': 86,
	'ṿ': 86,
	'Ẁ': 87,
	'ẁ': 87,
	'Ẃ': 87,
	'ẃ': 87,
	'Ẅ': 87,
	'ẅ': 87,
	'Ẇ': 87,
	'ẇ': 87,
	'Ẉ': 87,
	'ẉ': 87,
	'Ẋ': 88,
	'ẋ': 88,
	'Ẍ': 88,
	'ẍ': 88,
	'Ẏ': 89,
	'ẏ': 89,
	'Ẑ': 90,
	'ẑ': 90,
	'Ẓ': 90,
	'ẓ': 90,
	'Ẕ': 90,
	'ẕ': 90,
	'ẖ': 72,
	'ẗ': 84,
	'ẘ': 87,
	'ẙ': 89,
	'Ạ': 65,
	'ạ': 65,
	'Ả': 65,
	'ả': 65,
	'Ấ': 65,
	'ấ': 65,
	'Ầ': 65,
	'ầ': 65,
	'Ẩ': 65,
	'ẩ': 65,
	'Ẫ': 65,
	'ẫ': 65,
	'Ậ': 65,
	'ậ': 65,
	'Ắ': 65,
	'ắ': 65,
	'Ằ': 65,
	'ằ': 65,
	'Ẳ': 65,
	'ẳ': 65,
	'Ẵ': 65,
	'ẵ': 65,
	'Ặ': 65,
	'ặ': 65,
	'Ẹ': 69,
	'ẹ': 69,
	'Ẻ': 69,
	'ẻ': 69,
	'Ẽ': 69,
	'ẽ': 69,
	'Ế': 69,
	'ế': 69,
	'Ề': 69,
	'ề': 69,
	'Ể': 69,
	'ể': 69,
	'Ễ': 69,
	'ễ': 69,
	'Ệ': 69,
	'ệ': 69,
	'Ỉ': 73,
	'ỉ': 73,
	'Ị': 73,
	'ị': 73,
	'Ọ': 79,
	'ọ': 79,
	'Ỏ': 79,
	'ỏ': 79,
	'Ố': 79,
	'ố': 79,
	'Ồ': 79,
	'ồ': 79,
	'Ổ': 79,
	'ổ': 79,
	'Ỗ': 79,
	'ỗ': 79,
	'Ộ': 79,
	'ộ': 79,
	'Ớ': 79,
	'ớ': 79,
	'Ờ': 79,
	'ờ': 79,
	'Ở': 79,
	'ở': 79,
	'Ỡ': 79,
	'ỡ': 79,
	'Ợ': 79,
	'ợ': 79,
	'Ụ': 85,
	'ụ': 85,
	'Ủ': 85,
	'ủ': 85,
	'Ứ': 85,
	'ứ': 85,
	'Ừ': 85,
	'ừ': 85,
	'Ử': 85,
	'ử': 85,
	'Ữ': 85,
	'ữ': 85,
	'Ự': 85,
	'ự': 85,
	'Ỳ': 89,
	'ỳ': 89,
	'Ỵ': 89,
	'ỵ': 89,
	'Ỷ': 89,
	'ỷ': 89,
	'Ỹ': 89,
	'ỹ': 89,
	'K': 75,
	'Å': 65,
}

var alpha = map[rune]bool{}

func init() {
	for key := range alphaCode {
		alpha[key] = true
	}
}

func New(lower bool) SimilText {
	// Sets the alphabet of the parent characters.
	// If lower is true, set to lowercase, otherwise, set to uppercase
	st := SimilText{parent: map[int32]rune{}}

	st.parent[70] = iifRune(lower, 'f', 'F')
	st.parent[71] = iifRune(lower, 'g', 'G')
	st.parent[72] = iifRune(lower, 'h', 'H')
	st.parent[74] = iifRune(lower, 'j', 'J')
	st.parent[75] = iifRune(lower, 'k', 'K')
	st.parent[76] = iifRune(lower, 'l', 'L')
	st.parent[77] = iifRune(lower, 'm', 'M')
	st.parent[81] = iifRune(lower, 'q', 'Q')
	st.parent[82] = iifRune(lower, 'r', 'R')
	st.parent[84] = iifRune(lower, 't', 'T')
	st.parent[86] = iifRune(lower, 'v', 'V')
	st.parent[87] = iifRune(lower, 'w', 'W')
	st.parent[88] = iifRune(lower, 'x', 'X')
	st.parent[65] = iifRune(lower, 'a', 'A')
	st.parent[69] = iifRune(lower, 'e', 'E')
	st.parent[73] = iifRune(lower, 'i', 'I')
	st.parent[79] = iifRune(lower, 'o', 'O')
	st.parent[85] = iifRune(lower, 'u', 'U')
	st.parent[67] = iifRune(lower, 'c', 'C')
	st.parent[78] = iifRune(lower, 'n', 'N')
	st.parent[83] = iifRune(lower, 's', 'S')
	st.parent[90] = iifRune(lower, 'z', 'Z')
	st.parent[89] = iifRune(lower, 'y', 'Y')
	st.parent[68] = iifRune(lower, 'd', 'D')
	st.parent[80] = iifRune(lower, 'p', 'P')
	st.parent[66] = iifRune(lower, 'b', 'B')

	return st
}

func iifRune(lower bool, trueRune, falseRune rune) rune {
	if lower {
		return trueRune
	} else {
		return falseRune
	}
}

// Devuelve el código del alfabeto de similitudes correspondiente a una runa.
// Si el carácter no está el el alfabeto se devuelve la propia runa.
func code(c rune) int32 {
	if p, ok := alphaCode[c]; ok {
		return p
	} else {
		return c
	}
}

func isAlphaNum(c rune) bool {
	return (c >= '0' && c <= '9') || c == '_' || alpha[c]
}

func WordsRaw(s string, noRepeat bool) []string {
	start, length := 0, 0
	list := make([]string, 0, 2)

	for i, r := range s {
		if isAlphaNum(r) {
			if length == 0 {
				start = i
			}
			length += utf8.RuneLen(r)
		} else if length > 0 {
			word := s[start : start+length]
			if !(noRepeat && slices.Contains(list, word)) {
				list = append(list, word)
			}
			length = 0
		}
	}

	if length > 0 {
		word := s[start : start+length]
		if !(noRepeat && slices.Contains(list, word)) {
			list = append(list, word)
		}
	}

	return list[:len(list):len(list)]
}

func NormalizeRaw(s string) string {
	return strings.Join(WordsRaw(s, false), " ")
}

func (st SimilText) ParentRune(c rune) rune {
	c2, ok := st.parent[code(c)]

	if ok {
		return c2
	} else {
		return c
	}
}

func (st SimilText) ParentString(s string) string {
	sb := strings.Builder{}

	for _, r := range s {
		sb.WriteRune(st.ParentRune(r))
	}

	return sb.String()
}

func (st SimilText) Compare(s1, s2 string) int {
	reader1 := strings.NewReader(s1)
	reader2 := strings.NewReader(s2)

	for {
		r1, _, err1 := reader1.ReadRune()
		r2, _, err2 := reader2.ReadRune()

		if err1 == io.EOF && err2 == io.EOF {
			return 0
		} else if err1 == io.EOF {
			return -1
		} else if err2 == io.EOF {
			return 1
		}
		i, j := code(r1), code(r2)
		if i > j {
			return 1
		} else if i < j {
			return -1
		}
	}
}

func (st SimilText) SimilLevenshtein(s, t string) float64 {
	m := utf8.RuneCountInString(s)
	n := utf8.RuneCountInString(t)

	// Verifica que exista algo que comparar
	if m == 0 && n == 0 {
		return 100
	} else if m == 0 || n == 0 {
		return 0
	}

	// Generamos espacio de almacenamiento, (m + 1) filas y (n + 1) columnas
	d := make([][]int, m+1)
	for i := 0; i < len(d); i++ {
		d[i] = make([]int, n+1)
	}

	// Llena la primera columna y la primera fila.
	for i := 0; i <= m; i++ {
		d[i][0] = i
	}
	for j := 0; j <= n; j++ {
		d[0][j] = j
	}

	// Recorre la matriz llenando cada unos de los pesos
	// i filas, j columnas
	i := 1
	for _, ri := range s {
		j := 1
		for _, rj := range t {
			// Si son iguales en posiciones equidistantes el peso es 0
			// de lo contrario el peso suma a uno.
			cost := 0
			if code(ri) != code(rj) {
				cost = 1
			}
			d[i][j] = min(min(
				// Eliminación
				d[i-1][j]+1,
				// Inserción
				d[i][j-1]+1),
				// Sustitución
				d[i-1][j-1]+cost)
			j++
		}
		i++
	}

	// Calculamos el porcentaje de cambios en la palabra
	if m > n {
		return 100 * (1 - float64(d[m][n])/float64(m))
	} else {
		return 100 * (1 - float64(d[m][n])/float64(n))
	}
}

func (st SimilText) Words(s string, parent, noRepeat bool) []string {
	if parent {
		s = st.ParentString(s)
	}
	return WordsRaw(s, noRepeat)
}

// Normaliza una cadena. Por ejemplo para una instancia {st SimilText}:
//
//	st.Normalize("hola ,  mundo", "-", true) => "HOLA-MUNDO"
//	st.Normalize("123, responda, otra vez", "_", false) => "123_responda_otra_vez"
func (st SimilText) Normalize(s, sep string, parent bool) string {
	return strings.Join(st.Words(s, parent, false), sep)
}

func (st SimilText) Simil(s1, s2 string, subset, penalizeNumWords bool) float64 {
	ws1 := st.Words(s1, true, subset)
	ws2 := st.Words(s2, true, subset)

	if len(ws1) == 0 && len(ws2) == 0 {
		return 100
	} else if len(ws1) == 0 || len(ws2) == 0 {
		return 0
	}

	if len(ws1) > len(ws2) {
		ws1, ws2 = ws2, ws1
	}

	total, p := 0.0, 0.0
	count := 0
	// Buscamos inicialmente coincidencias del 100%
	for i := 0; i < len(ws1); i++ {
		for j := 0; j < len(ws2); j++ {
			if ws2[j] == "" {
				continue
			}
			if st.Compare(ws1[i], ws2[j]) != 0 {
				continue
			}
			pen := 1.0
			if i != j && !subset {
				pen = similFactorPen
			}
			total += 100 * pen
			count++
			ws1[i] = ""
			ws2[j] = ""
			break
		}
	}

	for count < len(ws1) {
		max := 0.0
		iMax, jMax := -1, -1
		for i := 0; i < len(ws1); i++ {
			if ws1[i] == "" {
				continue
			}
			for j := 0; j < len(ws2); j++ {
				if ws2[j] == "" {
					continue
				}
				p = st.SimilLevenshtein(ws1[i], ws2[j])
				// Se penaliza si la coincidencia es menor del mínimo
				if p < similMinPercent {
					p *= similFactorPen
				}
				// Se aplica penalización si están en distinta
				// posición y no se busca coincidencia de subconjunto
				if i != j && !subset {
					p *= similFactorPen
				}

				if p > max {
					max = p
					iMax = i
					jMax = j
					// Si coincide totalmente, no se sigue comprobando
					if max == 100 {
						break
					}
				}
			}

			// Si coincide totalmente, no se sigue comprobando
			if max == 100 {
				break
			}
		}

		if max > 0 {
			total += max
			ws1[iMax] = ""
			ws2[jMax] = ""
			count++
		} else {
			break
		}
	}

	// Si se compara como subconjunto, entonces se usa la cadena más corta
	if subset {
		p = total / float64(len(ws1))
	} else {
		p = total / float64(len(ws2))
	}
	if penalizeNumWords {
		p *= math.Pow(similNumWordsFactorPen, math.Abs(float64(len(ws1)-len(ws2))))
	}
	return p
}

// Search the 'search' string  in 'str' string. If found then return
// then 'start' index of first byte and the 'size' size in bytes.
//
// If not found then return 'start' as -1 value and 'size' equal to 0.
func (st SimilText) IndexOf(str, search string) (start int, size int) {
	if len(search) == 0 {
		return 0, 0
	} else if len(str) == 0 {
		return -1, 0
	}

	readerStr := strings.NewReader(str)
	readerSearch := strings.NewReader(search)

	// start, size = 0, 0 // by default
	for {
		strRune, strSize, strErr := readerStr.ReadRune()
		searchRune, _, searchErr := readerSearch.ReadRune()

		if searchErr == io.EOF {
			return start, size
		} else if strErr == io.EOF {
			return -1, 0
		}

		i, j := code(strRune), code(searchRune)
		if i == j {
			size += strSize
		} else {
			// reset search from next rune to current start index
			readerStr.Seek(-int64(size+strSize), io.SeekCurrent)
			strRune, strSize, strErr = readerStr.ReadRune()
			start += strSize
			readerSearch.Seek(0, io.SeekStart)
			size = 0
		}
	}
}
