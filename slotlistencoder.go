package main

import (
	"bytes"
	"github.com/blang/e12bot/parsing"
	"strconv"
)

const errorStr = "Ich konnte bisher keine Slotliste bekommen, aber vll kommt da ja noch was!\n"

func EncodeSlotList(slotlist *parsing.SlotList) string {
	if slotlist == nil {
		return errorStr
	}
	buff := bytes.NewBufferString("")
	buff.WriteString("# Slotliste\n")
	buff.WriteString("*(Update alle 5 Minuten)*\n\n")
	if len(slotlist.SlotListGroups) > 20 {
		return errorStr
	}
	for _, g := range slotlist.SlotListGroups {
		if len(g.Slots) > 20 {
			return errorStr
		}
		if g.Name != "" {
			buff.WriteString("##" + g.Name + "\n")
		}
		if g.Description != "" {
			buff.WriteString(g.Description + "\n")
		}
		for _, u := range g.Slots {
			if u.Number > 0 {
				buff.WriteString(strconv.Itoa(u.Number) + " ")
			}
			buff.WriteString(u.Name + ": ")
			if u.User != "" {
				buff.WriteString(u.User)
			} else {
				buff.WriteString("Frei")
			}
			buff.WriteString("\n\n")
		}
	}

	return buff.String()
}
