package main

import (
	"bytes"
	p "github.com/blang/e12bot/parsing"
	"html/template"
	"testing"
)

const tmpl = `
{{with .SlotListGroups}}
    {{range .}}
        Groupname {{.Name}}
        Description {{.Description}}
        {{range .Slots}}
        	Nr {{.Number}}
					Name {{.Name}}
					User {{.Name}}
        {{end}}
    {{end}}
{{end}}

`

var slotlist = &p.SlotList{
	SlotListGroups: []*p.SlotListGroup{
		&p.SlotListGroup{
			Name:        "Gruppe1",
			Description: "Beschreibung",
			Slots: []*p.SlotListSlot{
				&p.SlotListSlot{
					Number: 1,
					Name:   "Slot",
					User:   "Test",
				},
			},
		},
	},
}

func TestTemplate(t *testing.T) {
	tmpl, err := template.New("foo").Parse(tmpl)
	if err != nil {
		t.Errorf("Error while new template: %s", err)
	}
	buff := bytes.NewBufferString("")
	err = tmpl.Execute(buff, slotlist)
	if err != nil {
		t.Errorf("Error while new template: %s", err)
	}
	t.Logf("Output: %s", buff.String())
}
