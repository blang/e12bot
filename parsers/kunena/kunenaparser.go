package kunena

import (
	"github.com/blang/e12bot/parsing"
	"regexp"
	"strconv"
	"strings"
)

var slotlistBeginRegex = regexp.MustCompile("([Ss]lotlist|[Ss]lots|[Ss]lotliste)")
var patternSlot = regexp.MustCompile("#*\\s*[^\\W\\d]*([\\d]+)\\s*([\\w\\s-\\d]+)\\s*:\\s*(.*?)$")
var patternTag = regexp.MustCompile("(\\<.*?\\>)")
var patternGroup = regexp.MustCompile("<b>(.*?)</b>.*?<br.*?>")
var slotlistEndRegex = regexp.MustCompile("</td>")

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
	group := &parsing.SlotListGroup{}
	inSlotlist := false
	for _, line := range lines {
		if !inSlotlist {
			if slotlistBeginRegex.MatchString(line) {
				inSlotlist = true
				continue
			}
		} else {
			if slotlistEndRegex.MatchString(line) {
				break
			} else if patternSlot.MatchString(line) {
				parseSlot(line, group)
			} else if patternGroup.MatchString(line) {
				group = parseGroup(line, slotlist, group)
			}
		}
	}
	if len(group.Slots) > 0 {
		slotlist.SlotListGroups = append(slotlist.SlotListGroups, group)
	}
	return slotlist
}

func parseSlot(line string, g *parsing.SlotListGroup) {
	m := patternSlot.FindStringSubmatch(line)
	slot := &parsing.SlotListSlot{}
	if len(m) > 1 {
		num, err := strconv.Atoi(m[1])
		if err == nil {
			slot.Number = num
		}
	}
	if len(m) > 2 {
		slot.Name = sanitize(m[2])
	}
	if len(m) > 3 {
		slot.User = sanitize(m[3])
	}
	g.Slots = append(g.Slots, slot)
}

func parseGroup(line string, slotlist *parsing.SlotList, group *parsing.SlotListGroup) *parsing.SlotListGroup {
	m := patternGroup.FindStringSubmatch(line)

	var name string
	if len(m) > 1 {
		name = sanitize(m[1])
	}
	if name != "" {
		if len(group.Slots) > 0 {
			slotlist.SlotListGroups = append(slotlist.SlotListGroups, group)
		}
		slotgroup := &parsing.SlotListGroup{}
		slotgroup.Name = name
		return slotgroup
	}

	return group
}

func sanitize(text string) string {
	for patternTag.MatchString(text) {
		m := patternTag.FindStringSubmatch(text)
		text = strings.Replace(text, m[1], "", -1)
	}
	text = strings.Trim(text, "\n\r\t ")
	return text
}
