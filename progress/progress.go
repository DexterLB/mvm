package progress

// Progress
type Progress interface {
	SetTotal(total int)
	Set(value int)
	AddTotal(toTotal int)
	Add(toValue int)
	Done()
}

// ProgressChannel provides a way to communicate progress to a caller
type ProgressChannel struct {
	total        chan int
	value        chan int
	done         chan struct{}
	currentTotal int
	currentValue int
}

// NewProgressChannel initializes a progress channel object
func NewProgressChannel() *ProgressChannel {
	return &ProgressChannel{
		total: make(chan int),
		value: make(chan int),
	}
}

// DummyProgress doesn't report progress to anyone
func DummyProgress() Progress {
	p := NewProgressChannel()
	go func(p *ProgressChannel) {
		for _ = range p.Total() {
		}
	}(p)
	go func(p *ProgressChannel) {
		for _ = range p.Value() {
		}
	}(p)
	return p
}

func (p *ProgressChannel) SetTotal(total int) {
	p.currentTotal = total
	p.total <- total
}

func (p *ProgressChannel) Set(value int) {
	p.currentValue = value
	p.value <- value
}

func (p *ProgressChannel) AddTotal(toTotal int) {
	p.currentTotal += toTotal
	p.total <- p.currentTotal
}

func (p *ProgressChannel) Add(value int) {
	p.currentValue += value
	p.value <- p.currentValue
}

func (p *ProgressChannel) Total() <-chan int {
	return p.total
}

func (p *ProgressChannel) Done() {
	close(p.value)
	close(p.total)
	close(p.done)
}

func (p *ProgressChannel) Value() <-chan int {
	return p.value
}

func (p *ProgressChannel) DoneChan() <-chan struct{} {
	return p.done
}
