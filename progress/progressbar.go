package progress

import "github.com/cheggaaa/pb"

// ProgressBar is a console progress bar which displays much info
type ProgressBar struct {
	pb.ProgressBar
}

// NewProgressBar creates a progressbar with a 0/total progress. It's pretty
// useless until you call Start() to make it show up.
func NewProgressBar(total int) *ProgressBar {
	p := &ProgressBar{
		ProgressBar: *pb.New(total),
	}
	p.ShowPercent = true
	p.ShowTimeLeft = true
	return p
}

// StartProgressBar creates a progressbar and makes it show up immediately.
// Prefix shows as a title to the progressbar.
func StartProgressBar(total int, prefix string) *ProgressBar {
	p := NewProgressBar(total)
	p.Prefix(prefix)
	p.Start()
	return p
}

func (p *ProgressBar) SetTotal(max int) {
	p.Total = int64(max)
}

func (p *ProgressBar) AddTotal(toTotal int) {
	p.Total += int64(toTotal)
}

func (p *ProgressBar) Done() {
	p.Finish()
}

func (p *ProgressBar) Add(value int) {
	p.ProgressBar.Add(value)
}

func (p *ProgressBar) Set(value int) {
	p.ProgressBar.Set(value)
}
