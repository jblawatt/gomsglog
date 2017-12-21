package main

import (
	"github.com/jblawatt/gomsglog/gomsglog/parsers"
)

var p *parsers.TagsParser

func Parse(m *parsers.Message, rm *parsers.ReplacementManager) error {
	return p.Parse(m, rm)
}
