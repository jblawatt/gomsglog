package parsers

import (
	"fmt"
	"regexp"
	"strconv"
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

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

const ATTR_FULL = 1
const ATTR_KEY = 2
const ATTR_VALUE = 3

func (p AttrsParser) Parse(m *Message, rm *ReplacementManager) error {
	attrRe := regexp.MustCompile(`\$((\w+)[=:]([\w\-]+))`)
	for _, match := range attrRe.FindAllStringSubmatch(m.Original, -1) {
		replID := NewUUID()
		m.HTML = strings.Replace(m.HTML, match[0], replID, 1)
		rm.Add(
			replID,
			fmt.Sprintf(
				`<span class="gml-item gml-item__attr" data-gml-key="%s" data-gml-value="%s">%s</span>`,
				match[ATTR_KEY],
				match[ATTR_VALUE],
				match[ATTR_FULL],
			),
		)

		parserHit := false

		for _, p := range attrParser {
			foundSlug := strings.ToLower(match[2])
			if p.GetSlug() == foundSlug || contains(p.GetAlias(), foundSlug) {
				m.Attributes = append(m.Attributes, p.Parse(match[ATTR_VALUE]))
				parserHit = true
				break
			}
		}

		if !parserHit {
			m.Attributes = append(m.Attributes, Attr{
				Slug:        match[ATTR_KEY],
				ScreenName:  match[ATTR_KEY],
				Type:        "string",
				StringValue: match[ATTR_VALUE],
			})
		}

	}
	return nil
}

type AttrParser interface {
	GetSlug() string
	GetAlias() []string
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

func (d *DueAttrParser) GetAlias() []string {
	return []string{}
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

type HoursAttrParser struct{}

func (h *HoursAttrParser) GetSlug() string {
	return "hours"
}

func (h *HoursAttrParser) GetAlias() []string {
	return []string{"h"}
}

func (h *HoursAttrParser) Parse(value string) Attr {
	valueInt, _ := strconv.Atoi(value)
	return Attr{
		Slug:        "hours",
		ScreenName:  "hours",
		Type:        "int",
		StringValue: value,
		IntValue:    int64(valueInt),
	}
}

type QuoteAttrParser struct{}

func (q *QuoteAttrParser) GetSlug() string {
	return "quote"
}

func (h *QuoteAttrParser) GetAlias() []string {
	return []string{"quo"}
}

func (h *QuoteAttrParser) Parse(value string) Attr {
	valueInt, _ := strconv.Atoi(value)
	return Attr{
		Slug:        "quote",
		ScreenName:  "quote",
		Type:        "int",
		StringValue: value,
		IntValue:    int64(valueInt),
	}
}

func init() {
	RegisterParser(&AttrsParser{})
	RegisterAttrParser(&DueAttrParser{})
	RegisterAttrParser(&HoursAttrParser{})
	RegisterAttrParser(&QuoteAttrParser{})
}
