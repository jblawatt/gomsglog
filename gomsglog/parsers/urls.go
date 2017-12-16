package parsers

import (
	"fmt"
	"regexp"
	"strings"
)

type URLsParser struct{}

func (u URLsParser) Parse(m *Message, rm *ReplacementManager) {
	urlsRe := regexp.MustCompile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`)
	for _, match := range urlsRe.FindAllStringSubmatch(m.Original, -1) {
		replID := NewUUID()
		m.HTML = strings.Replace(m.HTML, match[0], replID, 1)
		rm.Add(
			replID,
			fmt.Sprintf(
				`<a class="gml-item gml-item__link" href="%s">%s</a>`,
				match[0],
				match[0],
			),
		)
		m.URLs = append(m.URLs, match[0])
	}
}

func init() {
	RegisterParser(&URLsParser{})
}
