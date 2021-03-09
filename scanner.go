// Package strscan is aiming for usability of Ruby StringScanner in Go.
package strscan

import (
	"regexp"

	"golang.org/x/exp/utf8string"
)

type location struct {
	lineno int
	column int
}

// StringScanner is a object to scan string.
//
// In little more detail, it contains a subject string and
// 'scan pointer', an index number to point to where in a subject string.
// User trys to match a head of pointed substring. If it matched,
// 'scan pointer' were progressed to next of matched part.
//
type StringScanner struct {
	subject     *utf8string.String
	rest        string
	pos         int
	lastPos     int
	matched     string
	lastMatched bool
	locMemo     map[int]*location
}

// New is a constructor for StringScanner.
func New(subject string) *StringScanner {
	return &StringScanner{
		subject: utf8string.NewString(subject),
		pos:     0,
		lastPos: -1,
		locMemo: map[int]*location{},
	}
}

// IsEOF is returns true if all of subject string were scanned.
func (s *StringScanner) IsEOF() bool { return s.pos >= s.subject.RuneCount() }

// Pos is returns 'scan pointer' value.
func (s *StringScanner) Pos() int { return s.pos }

// SetPos is move 'scan pointer'.
func (s *StringScanner) SetPos(n int) { s.pos = n }

// Matched is returns string which is a latest macthed part.
func (s *StringScanner) Matched() string { return s.matched }

// IsMatched is returns true if a latest matching was succeeded.
func (s *StringScanner) IsMatched() bool { return s.lastMatched }

func (s *StringScanner) updateRest() {
	if s.pos == s.lastPos {
		return
	}
	s.lastPos = s.pos
	s.rest = s.subject.Slice(s.pos, s.subject.RuneCount())
}

// Scan trys to match re to a head of pointed substring and
// when it's succeeded, progresses 'scan pointer' to next of
// matched part.
func (s *StringScanner) Scan(re *regexp.Regexp) bool {
	s.updateRest()
	loc := re.FindStringIndex(s.rest)
	if loc == nil || loc[0] != 0 {
		s.lastMatched = false
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
	part := utf8string.NewString(s.subject.Slice(0, pos))
	lineno := 0
	column := 0
	for i := 0; i < part.RuneCount(); i++ {
		if part.At(i) == '\n' {
			lineno++
			column = 0
		}
		column++
	}
	s.locMemo[pos] = &location{
		lineno: lineno,
		column: column,
	}
	return lineno, column
}
