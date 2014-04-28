package parsing

type SlotListParser interface {
	Parse(input string, url string) *SlotList
	Accept(url string) bool
}

type ParserCollection struct {
	parsers []SlotListParser
}

func (c *ParserCollection) Handle(p SlotListParser) {
	c.parsers = append(c.parsers, p)
}

func (c *ParserCollection) Parse(input string, url string) *SlotList {
	var hiscore float32
	hiscore = 0
	var hisl *SlotList

	for _, p := range c.parsers {
		if !p.Accept(url) {
			continue
		}
		sl := p.Parse(input, url)
		if sl != nil {
			if cscore := Classify(sl); cscore > hiscore {
				hisl = sl
			}
		}
	}
	return hisl
}

func (c *ParserCollection) Accept(url string) bool {
	for _, p := range c.parsers {
		if p.Accept(url) {
			return true
		}
	}
	return false
}
