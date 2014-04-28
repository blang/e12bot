package wiki

import (
	"testing"
)

const validSlotList = `
{| class="wikitable"
|-
! Nummer  !! Slot !! Besetzung
|-
| colspan="3"| '''LEAD'''
|-
| 1 || Squadlead|| '''Kerodan'''
|-
| 2 || Squadmedic|| '''Obi'''
|-
| 3 || Forward Air Commander || '''Spirit'''
|-
| colspan="3"| FAC Callsign: '''WIZARD'''
|-
| colspan="3" |
|-
| colspan="3"| '''ALPHA'''
|-
| 4 || Fireteam-Leader|| '''Ugene'''
|-
| 5 || Automatic Rifleman|| '''Systemstoerung81'''
|-
| 6 || AR Assistant|| '''Coati'''
|-
| 7 || Grenadier|| '''Paul'''
|-
|}
`

func TestParseWiki(t *testing.T) {
	w := &WikiParser{}
	sl := w.Parse(validSlotList, "http://wiki.echo12.de/wiki/Example?action=raw")
	if sl == nil {
		t.Error("Can't parse slot list")
	}

}
