package main

import (
	"fmt"
	"regexp"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func newUUID() string {
	return uuid.NewV4().String()
}

func HandleUsers(message *Message, replacements *map[string]string) {
	userRe := regexp.MustCompile(`@(?P<user>\w+)`)
	for _, match := range userRe.FindAllStringSubmatch(message.Original, -1) {
		replID := newUUID()
		message.HTML = strings.Replace(message.HTML, match[0], replID, 1)
		(*replacements)[replID] = fmt.Sprintf(
			`<span class="gml-item gml-item__user" data-user="%s">%s</span>`,
			match[0],
			match[1],
		)
		message.RelatedUsers = append(message.RelatedUsers, strings.ToLower(match[1]))
	}
}

func HandleTags(message *Message, replacements *map[string]string) {
	tagsRe := regexp.MustCompile(`#(\w+)`)
	for _, match := range tagsRe.FindAllStringSubmatch(message.Original, -1) {
		replID := newUUID()
		message.HTML = strings.Replace(message.HTML, match[0], replID, 1)
		(*replacements)[replID] = fmt.Sprintf(
			`<span class="gml-item gml-item__tag" data-tag="%s">%s</span>`,
			match[0],
			match[1],
		)
		message.Tags = append(message.Tags, strings.ToLower(match[1]))
	}
}

func HandleLinks(message *Message, replacements *map[string]string) {
	urlsRe := regexp.MustCompile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`)
	for _, match := range urlsRe.FindAllStringSubmatch(message.Original, -1) {
		replID := newUUID()
		message.HTML = strings.Replace(message.HTML, match[0], replID, 1)
		(*replacements)[replID] = fmt.Sprintf(
			`<a class="gml-item gml-item__link" href="%s">%s</span>`,
			match[0],
			match[0],
		)
		message.URLs = append(message.URLs, match[0])
	}
}

func HandleAttrs(message *Message, replacements *map[string]string) {
	attrRe := regexp.MustCompile(`\$(\w+):([\w\-]+)`)
	for _, match := range attrRe.FindAllStringSubmatch(message.Original, -1) {
		replID := newUUID()
		message.HTML = strings.Replace(message.HTML, match[0], replID, 1)
		(*replacements)[replID] = fmt.Sprintf(
			`<span class="gml-item gml-item__attr" data-gml-key="%s" data-gml-value="%s">%s</span>`,
			match[1],
			match[2],
			match[2],
		)
		message.Attributes[match[1]] = match[2]
	}

}
