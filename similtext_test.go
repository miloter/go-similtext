package similtext

import (
	"testing"
)

var st SimilText

func setup() {
	st = New(false)
}

func teardown() {

}

func TestMain(m *testing.M) {
	setup()
	defer teardown()

	m.Run()
}

func TestParentString(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"áèîöuñ", "AEIOUN"},
		{"€uroç", "€UROC"},
		{"žůžo", "ZUZO"},
	}

	for _, ts := range tests {
		out := st.ParentString(ts.input)
		if out != ts.output {
			t.Errorf("Dado '%s', se esperaba '%s', pero se obtuvo '%s'",
				ts.input, ts.output, out)
		}
	}
}

func TestICmp(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		expected int
	}{
		{"holá", "hola", 0},
		{"hola", "HOLÁ", 0},
		{"hÒla mundo", "HöLA", 1},
		{"hola mundo y", "HOLA MUNDO Z", -1},
		{"žůžo", "ZUZO", 0},
		{"", "", 0},
	}

	for _, ts := range tests {
		out := st.Compare(ts.s1, ts.s2)
		if out != ts.expected {
			t.Errorf("Dadas %q, y %q se esperaba %d, pero se obtuvo %d",
				ts.s1, ts.s2, ts.expected, out)
		}
	}
}

func BenchmarkICmp(b *testing.B) {
	tests := []struct {
		s1       string
		s2       string
		expected int
	}{
		{"holá", "hola", 0},
		{"hola", "HOLÁ", 0},
		{"hÒla mundo", "HöLA", 1},
		{"hola mundo y", "HOLA MUNDO Z", -1},
		{"žůžo", "ZUZO", 0},
		{"", "", 0},
	}

	for _, ts := range tests {
		st.Compare(ts.s1, ts.s2)
	}
}
