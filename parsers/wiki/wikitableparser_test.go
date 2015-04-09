package wiki

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const validSlotList = `
{| class="wikitable"
|-
! Nummer  !! Slot !! Besetzung
|-
| colspan="3"| '''Gruppe1'''
|-
| 1 || Slot1|| '''User1'''
|-
| 2 || Slot2 ||
|-
| 3 ||Slot3|| '''User3'''
|-
| colspan="3"| Gruppe2
|-
|4 || Slot4|| '''User4'''
|-
|5||  Slot5  ||User5
|-
|  6  || Slot6|| '''User6'''
|-
| 7 || Slot7|| '''User7'''
|-
|}
`

const validSlotListDesc = `
{| class="wikitable"
|-
! Nummer  !! Slot !! Besetzung !! Bemerkung
|-
| colspan="4"| '''Gruppe1'''
|-
| 1 || Slot1|| '''User1''' || Test
|-
| 2 || Slot2 || ||
|-
| 3 ||Slot3|| '''User3''' ||
|-
| colspan="4"| Gruppe2
|-
|4 || Slot4|| '''User4''' ||
|-
|5||  Slot5  ||User5 ||
|-
|  6  || Slot6|| '''User6'''|| DLC-Slot
|-
| 7 || Slot7|| '''User7'''||
|-
|}
`

func TestParseWikiTable(t *testing.T) {
	Convey("Given a fresh parser", t, func() {

		w := &WikiTableParser{}
		Convey("When parsing wiki slotlist", func() {
			sl := w.Parse(validSlotList, "http://wiki.echo12.de/wiki/Example?action=raw")
			Convey("Slotlist can be parsed", func() {
				So(sl, ShouldNotBeNil)
			})
			Convey("Slotlist has proper amount of groups", func() {
				So(len(sl.SlotListGroups), ShouldEqual, 2)
				// group := sl.SlotListGroups[0]
			})
			Convey("First group", func() {
				g := sl.SlotListGroups[0]
				Convey("Has correct attributes", func() {
					So(g.Name, ShouldEqual, "Gruppe1")
					So(g.Description, ShouldEqual, "")
				})
				Convey("Slots", func() {
					Convey("Has 3 entries", func() {
						So(len(g.Slots), ShouldEqual, 3)
					})
					Convey("Slots are correct", func() {
						slot := g.Slots[0]
						So(slot.Name, ShouldEqual, "Slot1")
						So(slot.Number, ShouldEqual, 1)
						So(slot.User, ShouldEqual, "User1")
						slot = g.Slots[1]
						So(slot.Name, ShouldEqual, "Slot2")
						So(slot.Number, ShouldEqual, 2)
						So(slot.User, ShouldEqual, "")
						slot = g.Slots[2]
						So(slot.Name, ShouldEqual, "Slot3")
						So(slot.Number, ShouldEqual, 3)
						So(slot.User, ShouldEqual, "User3")
					})
				})
			})

			Convey("Second group", func() {
				g := sl.SlotListGroups[1]
				Convey("Has correct attributes", func() {
					So(g.Name, ShouldEqual, "Gruppe2")
					So(g.Description, ShouldEqual, "")
				})
				Convey("Slots", func() {
					Convey("Has 3 entries", func() {
						So(len(g.Slots), ShouldEqual, 4)
					})
					Convey("Slots are correct", func() {
						slot := g.Slots[0]
						So(slot.Name, ShouldEqual, "Slot4")
						So(slot.Number, ShouldEqual, 4)
						So(slot.User, ShouldEqual, "User4")
						slot = g.Slots[1]
						So(slot.Name, ShouldEqual, "Slot5")
						So(slot.Number, ShouldEqual, 5)
						So(slot.User, ShouldEqual, "User5")
						slot = g.Slots[2]
						So(slot.Name, ShouldEqual, "Slot6")
						So(slot.Number, ShouldEqual, 6)
						So(slot.User, ShouldEqual, "User6")
						slot = g.Slots[3]
						So(slot.Name, ShouldEqual, "Slot7")
						So(slot.Number, ShouldEqual, 7)
						So(slot.User, ShouldEqual, "User7")
					})
				})
			})
		})
	})
}

func TestParseWikiTableWithDescription(t *testing.T) {
	Convey("Given a fresh parser", t, func() {

		w := &WikiTableParser{}
		Convey("When parsing wiki slotlist", func() {
			sl := w.Parse(validSlotListDesc, "http://wiki.echo12.de/wiki/Example?action=raw")
			Convey("Slotlist can be parsed", func() {
				So(sl, ShouldNotBeNil)
			})
			Convey("Slotlist has proper amount of groups", func() {
				So(len(sl.SlotListGroups), ShouldEqual, 2)
				// group := sl.SlotListGroups[0]
			})
			Convey("First group", func() {
				g := sl.SlotListGroups[0]
				Convey("Has correct attributes", func() {
					So(g.Name, ShouldEqual, "Gruppe1")
					So(g.Description, ShouldEqual, "")
				})
				Convey("Slots", func() {
					Convey("Has 3 entries", func() {
						So(len(g.Slots), ShouldEqual, 3)
					})
					Convey("Slots are correct", func() {
						slot := g.Slots[0]
						So(slot.Name, ShouldEqual, "Slot1")
						So(slot.Number, ShouldEqual, 1)
						So(slot.User, ShouldEqual, "User1")
						So(slot.Desc, ShouldEqual, "Test")
						slot = g.Slots[1]
						So(slot.Name, ShouldEqual, "Slot2")
						So(slot.Number, ShouldEqual, 2)
						So(slot.User, ShouldEqual, "")
						slot = g.Slots[2]
						So(slot.Name, ShouldEqual, "Slot3")
						So(slot.Number, ShouldEqual, 3)
						So(slot.User, ShouldEqual, "User3")
					})
				})
			})

			Convey("Second group", func() {
				g := sl.SlotListGroups[1]
				Convey("Has correct attributes", func() {
					So(g.Name, ShouldEqual, "Gruppe2")
					So(g.Description, ShouldEqual, "")
				})
				Convey("Slots", func() {
					Convey("Has 3 entries", func() {
						So(len(g.Slots), ShouldEqual, 4)
					})
					Convey("Slots are correct", func() {
						slot := g.Slots[0]
						So(slot.Name, ShouldEqual, "Slot4")
						So(slot.Number, ShouldEqual, 4)
						So(slot.User, ShouldEqual, "User4")
						slot = g.Slots[1]
						So(slot.Name, ShouldEqual, "Slot5")
						So(slot.Number, ShouldEqual, 5)
						So(slot.User, ShouldEqual, "User5")
						slot = g.Slots[2]
						So(slot.Name, ShouldEqual, "Slot6")
						So(slot.Number, ShouldEqual, 6)
						So(slot.User, ShouldEqual, "User6")
						So(slot.Desc, ShouldEqual, "DLC-Slot")
						slot = g.Slots[3]
						So(slot.Name, ShouldEqual, "Slot7")
						So(slot.Number, ShouldEqual, 7)
						So(slot.User, ShouldEqual, "User7")
					})
				})
			})
		})
	})
}
