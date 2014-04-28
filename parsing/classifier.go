package parsing

func Classify(l *SlotList) float32 {
	if l == nil {
		return 0
	}
	return l.classify()
}

func (l *SlotList) classify() float32 {
	var class float32 = 1.0
	for _, g := range l.SlotListGroups {
		if g != nil {
			class += g.classify()
		}
	}
	return class
}

func (l *SlotListGroup) classify() float32 {
	var i float32 = 0
	if l.Name != "" {
		i += 1
	}
	if l.Description != "" {
		i += 1
	}
	f := i / 2
	for _, s := range l.Slots {
		f += s.classify()
	}
	return f
}

func (l *SlotListSlot) classify() float32 {
	var i float32 = 0
	if l.Number > 0 && l.Number < 99 {
		i += 1
	}
	if l.Name != "" {
		i += 1
	}
	if l.User != "" {
		i += 1
	}
	return i / 3
}
