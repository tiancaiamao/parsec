package parsec

type Parser interface {
	Parse(input string) (bool, string)
}

type notC struct {
	raw Parser
}

func NOT(p Parser) Parser {
	return notC{p}
}

func (c notC) Parse(input string) (bool, string) {
	p := Parser(c.raw)
	succ, _ := p.Parse(input)
	if succ {
		return false, input
	}
	return true, input
}

func OR(ps ...Parser) Parser {
	return orC(ps)
}

type orC []Parser

func (c orC) Parse(input string) (bool, string) {
	for _, p := range c {
		succ, remain := p.Parse(input)
		if succ {
			return succ, remain
		}
	}
	return false, input
}

func SEQ(ps ...Parser) Parser {
	return seqC(ps)
}

type seqC []Parser

func (c seqC) Parse(input string) (bool, string) {
	var succ bool
	remain := input
	for _, p := range c {
		succ, remain = p.Parse(remain)
		if !succ {
			return false, input
		}
		continue
	}
	return true, remain
}

type Literal string

func (s Literal) Parse(input string) (bool, string) {
	if len(input) >= len(s) && input[:len(s)] == string(s) {
		return true, input[len(s):]
	}
	return false, input
}

type Byte byte

func (s Byte) Parse(input string) (bool, string) {
	if len(input) > 0 && input[0] == byte(s) {
		return true, input[1:]
	}
	return false, input
}

func Range(start, end byte) Parser {
	return rangeLit{start, end}
}

type rangeLit struct {
	start byte
	end   byte
}

func (r rangeLit) Parse(input string) (bool, string) {
	if len(input) > 0 {
		if input[0] >= r.start && input[0] <= r.end {
			return true, input[1:]
		}
	}
	return false, input
}

type WhiteSpace struct{}

func (w WhiteSpace) Parse(input string) (bool, string) {
	var succ bool
	remain := input
	for len(remain) > 0 && remain[0] == ' ' {
		succ = true
		remain = remain[1:]
	}
	return succ, remain
}
