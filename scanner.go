// Package strscan is aiming for usability of Ruby StringScanner in Go.
package strscan

import (
	"fmt"
	"regexp"

	"golang.org/x/exp/utf8string"
)

type location struct {
	lineno int
	column int
}

// StringScanner is a object to scan string.
//
// In little more detail, it contains a target string and
// 'scan pointer', an index number to point to where in a target string.
// User trys to match a head of pointed substring. If it matched,
// 'scan pointer' were progressed to next of matched part.
//
type StringScanner struct {
	target      *utf8string.String
	rest        string
	pos         int
	lastPos     int
	matched     string
	lastMatched bool
	locMemo     map[int]*location
}

// New is a constructor for StringScanner.
func New(target string) *StringScanner {
	return &StringScanner{
		target:  utf8string.NewString(target),
		pos:     0,
		lastPos: -1,
		locMemo: map[int]*location{},
	}
}

// IsEOF is returns true if all of target string were scanned.
func (s *StringScanner) IsEOF() bool { return s.pos >= s.target.RuneCount() }

// Pos is returns 'scan pointer' value.
func (s *StringScanner) Pos() int { return s.pos }

// SetPos is move 'scan pointer'.
//
// When n is negative, it means an offset from end of target string.
//
// It returns error against to invalid index for target string.
//
func (s *StringScanner) SetPos(n int) error {
	p := n
	if p < 0 {
		p = s.target.RuneCount() + n
	}
	if p < 0 || p >= s.target.RuneCount() {
		return fmt.Errorf("index out of range: %d", n)
	}
	s.pos = p
	return nil
}

// Matched is returns string which is a latest macthed part.
func (s *StringScanner) Matched() string { return s.matched }

// IsMatched is returns true if a latest matching was succeeded.
func (s *StringScanner) IsMatched() bool { return s.lastMatched }

func (s *StringScanner) updateRest() {
	if s.pos == s.lastPos {
		return
	}
	s.lastPos = s.pos
	s.rest = s.target.Slice(s.pos, s.target.RuneCount())
}

func (s *StringScanner) resetLastMatched() {
	s.matched = ""
	s.lastMatched = false
}

// Scan trys to match re to a head of pointed substring and
// when it's succeeded, progresses 'scan pointer' to next of
// matched part.
func (s *StringScanner) Scan(re *regexp.Regexp) bool {
	s.updateRest()
	s.resetLastMatched()
	loc := re.FindStringIndex(s.rest)
	if loc == nil || loc[0] != 0 {
		return false
	}
	s.matched = s.rest[0:loc[1]]
	s.lastMatched = true
	s.pos += loc[1]
	return true
}

// LinenoAndColumn calculates line number and column index of
// location pos pointed.
func (s *StringScanner) LinenoAndColumn(pos int) (int, int) {
	if loc, ok := s.locMemo[pos]; ok {
		return loc.lineno, loc.column
	}
	part := utf8string.NewString(s.target.Slice(0, pos))
	lineno := 0
	column := 0
	for i := 0; i < part.RuneCount(); i++ {
		if part.At(i) == '\n' {
			lineno++
			column = 0
			continue
		}
		column++
	}
	s.locMemo[pos] = &location{
		lineno: lineno,
		column: column,
	}
	return lineno, column
}
