package parsers

import (
	"fmt"
	"regexp"
	"strings"
)

type TagsParser struct{}

func (t TagsParser) Parse(m *Message, rm *ReplacementManager) {
	tagsRe := regexp.MustCompile(`#(\w+)`)
	for _, match := range tagsRe.FindAllStringSubmatch(m.Original, -1) {
		replID := NewUUID()
		m.HTML = strings.Replace(m.HTML, match[0], replID, 1)
		rm.Add(
			replID,
			fmt.Sprintf(
				`<span class="gml-item gml-item__tag" data-tag="%s">%s</span>`,
				match[0],
				match[1],
			),
		)
		m.Tags = append(m.Tags, strings.ToLower(match[1]))
	}
}

func init() {
	RegisterParser(&TagsParser{})
}
