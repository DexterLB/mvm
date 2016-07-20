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
type DummyProgress struct{}

// SetTotal does nothing
func (d *DummyProgress) SetTotal(total int) {
}

// Set does nothing
func (p *DummyProgress) Set(value int) {
}

// AddTotal does nothing
func (p *DummyProgress) AddTotal(toTotal int) {
}

// Add does nothing
func (p *DummyProgress) Add(value int) {
}

// Done does nothing
func (p *DummyProgress) Done() {
}

// SetTotal sets the total progress
func (p *ProgressChannel) SetTotal(total int) {
	p.currentTotal = total
	p.total <- total
}

// Set updates the current progress
func (p *ProgressChannel) Set(value int) {
	p.currentValue = value
	p.value <- value
}

// AddTotal adds to the current total
func (p *ProgressChannel) AddTotal(toTotal int) {
	p.currentTotal += toTotal
	p.total <- p.currentTotal
}

// Add adds to the current value
func (p *ProgressChannel) Add(value int) {
	p.currentValue += value
	p.value <- p.currentValue
}

// Total returns a channel which receives updates to the total value
func (p *ProgressChannel) Total() <-chan int {
	return p.total
}

// Done closes the channels
func (p *ProgressChannel) Done() {
	close(p.value)
	close(p.total)
	close(p.done)
}

// Value returns a channel which receives updates to the current value
func (p *ProgressChannel) Value() <-chan int {
	return p.value
}

// DoneChan returns a channel which is closed when Done() is called
func (p *ProgressChannel) DoneChan() <-chan struct{} {
	return p.done
}
