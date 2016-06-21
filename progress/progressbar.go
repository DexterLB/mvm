package progress

import "github.com/cheggaaa/pb"

// ProgressBar is a console progress bar which displays much info
type ProgressBar struct {
	pb.ProgressBar
}

// NewProgressBar creates a progressbar with a 0/1 progress. It's pretty
// useless until you use SetTotal() to make the total something other than 1.
// Also, you have to call Start() to make it show up.
func NewProgressBar() *ProgressBar {
	return &ProgressBar{*pb.New(1)}
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
