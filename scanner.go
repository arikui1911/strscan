package strscan

import (
	"regexp"

	"golang.org/x/exp/utf8string"
)

type StringScanner struct {
	subject     *utf8string.String
	rest        string
	pos         int
	lastPos     int
	matched     string
	lastMatched bool
}

func New(subject string) *StringScanner {
	return &StringScanner{
		subject: utf8string.NewString(subject),
		pos:     0,
		lastPos: -1,
	}
}

func (s *StringScanner) IsEOF() bool { return s.pos >= s.subject.RuneCount() }

func (s *StringScanner) Pos() int { return s.pos }

func (s *StringScanner) SetPos(n int) { s.pos = n }

func (s *StringScanner) Matched() string { return s.matched }

func (s *StringScanner) IsMatched() bool { return s.lastMatched }

func (s *StringScanner) Scan(re *regexp.Regexp) bool {
	if s.pos != s.lastPos {
		s.lastPos = s.pos
		s.rest = s.subject.Slice(s.pos, s.subject.RuneCount())
	}
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
