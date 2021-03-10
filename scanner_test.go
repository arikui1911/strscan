package strscan_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/arikui1911/strscan"
)

func TestScanner(t *testing.T) {
	src := `
def hoge(n)
  puts "HOGE!" * n
end
`

	spaces := regexp.MustCompile(`[\s\n]+`)
	ident := regexp.MustCompile(`[a-zA-Z_][a-zA-Z_0-9]*`)
	num := regexp.MustCompile(`\d+`)
	str := regexp.MustCompile(`".*?"`)
	ch := regexp.MustCompile(`.`)

	s := strscan.New(src)

	for !s.IsEOF() {
		var msg string
		p := s.Pos()

		switch {
		case s.Scan(spaces):
			continue
		case s.Scan(ident):
			msg = fmt.Sprintf("IDENT: %s", s.Matched())
		case s.Scan(num):
			msg = fmt.Sprintf("NUMBER: %s", s.Matched())
		case s.Scan(str):
			msg = fmt.Sprintf("STRING: %s", s.Matched())
		case s.Scan(ch):
			msg = fmt.Sprintf("%s: %s", s.Matched(), s.Matched())
		default:
			panic("must not happen")
		}

		fl, fc := s.LinenoAndColumn(p)
		ll, lc := s.LinenoAndColumn(s.Pos() - 1)
		fmt.Printf("%d,%d,%d,%d: %s\n", fl, fc, ll, lc, msg)
	}
}
