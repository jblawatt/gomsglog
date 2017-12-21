package parsers

import (
	"fmt"
	"regexp"
	"strings"
)

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
		m.Attributes[match[1]] = match[2]
	}
	return nil
}

func init() {
	RegisterParser(&AttrsParser{})
}
