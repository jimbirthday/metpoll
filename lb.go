package metpoll

type LB interface {
	Cur() int
	Set(t int)
}

type RoundRobinBalance struct {
	total int
	cur   int
}

func (rrb *RoundRobinBalance) Set(t int) {
	rrb.total = t
}

func (rrb *RoundRobinBalance) Cur() int {
	i := rrb.cur
	if rrb.cur <= rrb.total {
		rrb.cur += 1
	} else {
		i = 1
		rrb.cur = 1
	}
	return i - 1
}
