package wiki

import (
	parser "github.com/blang/e12bot/parsing"
	"regexp"
	"strconv"
	"strings"
)

var slotRegex = regexp.MustCompile("^\\|\\s*(\\d+)\\s*\\|\\|\\s*([\\pL\\s-\\d]+).*?\\|\\|\\s*([\\pL\\s-\\d'*]+)?")
var slotDescRegex = regexp.MustCompile("^\\| .*?\\s+([\\pL\\s-\\d'*]+)")

type WikiTableParser struct {
}

func (w *WikiTableParser) Accept(url string) bool {
	return strings.Contains(url, "http://wiki.echo12.de/wiki/")
}

func (w *WikiTableParser) Parse(input string, url string) *parser.SlotList {
	if !strings.Contains(url, "http://wiki.echo12.de/wiki/") {
		return nil
	}

	lines := strings.Split(input, "\n")
	slotlist := &parser.SlotList{}
	group := &parser.SlotListGroup{}
	inSlotlist := false
	for _, line := range lines {
		if !inSlotlist && strings.HasPrefix(line, "{|") && strings.Contains(line, "wikitable") {
			inSlotlist = true
			continue
		}
		if inSlotlist && strings.HasPrefix(line, "|}") {
			inSlotlist = true
			break
		}

		if inSlotlist {
			if slotRegex.MatchString(line) {
				parseSlot(line, group)
			} else if slotDescRegex.MatchString(line) {
				group = parseGroup(line, slotlist, group)
			}
		}
	}
	if len(group.Slots) > 0 {
		slotlist.SlotListGroups = append(slotlist.SlotListGroups, group)
	}
	return slotlist

}

func sanitize(s string) string {
	s = strings.Replace(s, "'", "", -1)
	s = strings.Trim(s, " \t.:\r")
	return s
}

func parseSlot(line string, group *parser.SlotListGroup) {
	m := slotRegex.FindStringSubmatch(line)
	slot := &parser.SlotListSlot{}
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
	group.Slots = append(group.Slots, slot)
}

func parseGroup(line string, slotlist *parser.SlotList, group *parser.SlotListGroup) *parser.SlotListGroup {
	m := slotDescRegex.FindStringSubmatch(line)
	if len(group.Slots) > 0 {
		slotlist.SlotListGroups = append(slotlist.SlotListGroups, group)
	}
	slotgroup := &parser.SlotListGroup{}
	if len(m) > 1 {
		slotgroup.Name = sanitize(m[1])
	}
	return slotgroup
}
