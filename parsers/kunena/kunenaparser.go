package kunena

import (
	"github.com/blang/e12bot/parsing"
	"regexp"
	"strings"
)

var slotlistBeginRegex = regexp.MustCompile("(Slotlist|Slots|Slotliste)")

type KunenaParser struct {
}

func (p *KunenaParser) Accept(url string) bool {
	return strings.Contains(url, "heeresgruppe2012.de")
}

func (p *KunenaParser) Parse(input string, url string) *parsing.SlotList {
	if !strings.Contains(url, "heeresgruppe2012.de") {
		return nil
	}
	lines := strings.Split(input, "\n")
	slotlist := &parsing.SlotList{}
	// group := &parsing.SlotListGroup{}
	inSlotlist := false
	for _, line := range lines {
		if !inSlotlist {
			if slotlistBeginRegex.MatchString(line) {
				inSlotlist = true
				continue
			}
		}
	}
	return slotlist
}
