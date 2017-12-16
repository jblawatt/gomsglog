package parsers

import (
	"fmt"
	"regexp"
	"strings"
)

type UsersParser struct{}

func (u UsersParser) Parse(m *Message, rm *ReplacementManager) {
	userRe := regexp.MustCompile(`@(?P<user>\w+)`)
	for _, match := range userRe.FindAllStringSubmatch(m.Original, -1) {
		replID := NewUUID()
		m.HTML = strings.Replace(m.HTML, match[0], replID, 1)
		rm.Add(
			replID,
			fmt.Sprintf(
				`<span class="gml-item gml-item__user" data-user="%s">%s</span>`,
				match[0],
				match[1],
			),
		)
		m.RelatedUsers = append(m.RelatedUsers, strings.ToLower(match[1]))
	}
}

func init() {
	RegisterParser(&UsersParser{})
}
