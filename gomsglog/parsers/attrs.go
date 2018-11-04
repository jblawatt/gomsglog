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

func (p AttrsParser) Parse(m *Message, rm *ReplacementManager) error {
	attrRe := regexp.MustCompile(`\&(\w+)[=:]([\w\-]+)`)
	for _, match := range attrRe.FindAllStringSubmatch(m.Original, -1) {
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

func init() {
	RegisterParser(&AttrsParser{})
	RegisterAttrParser(&DueAttrParser{})
	RegisterAttrParser(&HoursAttrParser{})
}
