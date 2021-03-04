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
		p := s.Pos()
		switch {
		case s.Scan(spaces):
			// do nothing
		case s.Scan(ident):
			fmt.Printf("%d: IDENT: %s\n", p, s.Matched())
		case s.Scan(num):
			fmt.Printf("%d: NUMBER: %s\n", p, s.Matched())
		case s.Scan(str):
			fmt.Printf("%d: STRING: %s\n", p, s.Matched())
		case s.Scan(ch):
			fmt.Printf("%d: %s: %s\n", p, s.Matched(), s.Matched())
		default:
			panic("must not happen")
		}
	}
}
