package parsing

type SlotList struct {
	SlotListGroups []*SlotListGroup
}

type SlotListGroup struct {
	Name        string
	Description string
	Slots       []*SlotListSlot
}

type SlotListSlot struct {
	Number int
	Name   string
	User   string
}
