package parsers

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Attr struct {
	Slug        string
	ScreenName  string
	Type        string
	DateValue   time.Time
	IntValue    int64
	FloatValue  float64
	StringValue string
	BoolValue   bool
}

type AttrsParser struct{}

func (p AttrsParser) Parse(m *Message, rm *ReplacementManager) error {
	attrRe := regexp.MustCompile(`\&(\w+)[=:]([\w\-]+)`)
	for _, match := range attrRe.FindAllStringSubmatch(m.Original, -1) {
		fmt.Println("ATTR PARSER")
		replID := NewUUID()
		m.HTML = strings.Replace(m.HTML, match[0], replID, 1)
		rm.Add(
			replID,
			fmt.Sprintf(
				`<span class="gml-item gml-item__attr" data-gml-key="%s" data-gml-value="%s">%s</span>`,
				match[1],
				match[2],
				match[2],
			),
		)

		parserHit := false

		for _, p := range attrParser {
			fmt.Println("ATTR PARSER ", match[1])
			if p.GetSlug() == strings.ToLower(match[1]) {
				m.Attributes[match[1]] = p.Parse(match[2])
				fmt.Println(m.Attributes[match[1]])
				parserHit = true
				break
			}
		}

		if !parserHit {
			m.Attributes[match[1]] = Attr{
				Slug:        match[1],
				ScreenName:  match[1],
				Type:        "string",
				StringValue: match[2],
			}
		}

	}
	return nil
}

type AttrParser interface {
	GetSlug() string
	Parse(value string) Attr
}

var attrParser []AttrParser

func RegisterAttrParser(p AttrParser) {
	attrParser = append(attrParser, p)
}

type DueAttrParser struct{}

func (d *DueAttrParser) GetSlug() string {
	return "due"
}

func (d *DueAttrParser) Parse(value string) Attr {
	t, _ := time.Parse("2006-01-02", value)
	return Attr{
		Slug:        "due",
		ScreenName:  "due",
		Type:        "date",
		StringValue: value,
		DateValue:   t,
	}
}

func init() {
	RegisterParser(&AttrsParser{})
	RegisterAttrParser(&DueAttrParser{})
}
