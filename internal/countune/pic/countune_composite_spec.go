package pic

type countunePicSpec struct {
	barCount  int
	barColor1 string
	barColor2 string
	bgColor   string
}

type CountuneCompositeSpec struct {
	startNum  int
	totalBars int
	picSpecs  []countunePicSpec
}

func NewCountuneCompositeSpec(startNum int) *CountuneCompositeSpec {

	return &CountuneCompositeSpec{
		startNum:  startNum,
		totalBars: 0,
	}

}

func (s *CountuneCompositeSpec) AddPic(barCount int, bgColor, barColor1, barColor2 string) {

	s.picSpecs = append(s.picSpecs, countunePicSpec{
		barCount:  barCount,
		barColor1: barColor1,
		barColor2: barColor2,
		bgColor:   bgColor,
	})
	s.totalBars += barCount
}
