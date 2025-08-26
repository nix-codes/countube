package pic

import (
	"countube/internal/common"

	// "fmt"
	"image"
)

type Direction int

const (
	UP Direction = iota
	DOWN
)

type PicSeqSpec struct {
	PicHeightInPixels   int
	PixelsPerUnit       int
	AmplitudeUnits      int
	StartNum            int
	StartHeightUnits    int
	StartDirection      Direction
	InitialBarPlacement BarPlacement
}

type CountunePicSeqProvider struct {
	spec         PicSeqSpec
	nextNum      int
	barPlacement BarPlacement
}

func NewCountunePicProvider(spec PicSeqSpec) *CountunePicSeqProvider {

	// primesInRange := common.CountPrimesInRange(spec.StartNum, p.nextNum)
	// fmt.Printf("primes in range: %d\n", primesInRange)
	// barPlacement := p.spec.InitialBarPlacement

	// if primesInRange%2 == 1 {
	// 	barPlacement.Toggle()
	// }

	return &CountunePicSeqProvider{
		spec:         spec,
		nextNum:      spec.StartNum,
		barPlacement: spec.InitialBarPlacement,
	}
}

func (p *CountunePicSeqProvider) NextRandomPic() *image.RGBA {
	barCount := selectRandomPicSize()
	// barWidth := p.spec.PixelsPerUnit
	// picWidth := barCount * barWidth
	colors := selectRandomCountuneColors()
	// colors := []string{"0000ff", "ff0000", "00ff00"}

	// primesInRange := common.CountPrimesInRange(p.spec.StartNum, p.nextNum)
	// fmt.Printf("primes in range: %d\n", primesInRange)
	// barPlacement := p.spec.InitialBarPlacement

	// if primesInRange%2 == 1 {
	// 	barPlacement.Toggle()
	// }
	// fmt.Println("-")

	editor := NewPicEditor(p.spec.PicHeightInPixels, p.spec.PixelsPerUnit, barCount, p.spec.AmplitudeUnits, CENTER, p.barPlacement)
	editor.ChangeColors(hexToColor(colors[0]), hexToColor(colors[1]), hexToColor(colors[2]))

	// for n := p.nextNum; n < p.nextNum+barCount; n++ {

	// 	if common.IsPrime(n) {
	// 		p.barPlacement = p.barPlacement.Toggle()
	// 	}

	// 	y := countuneFn(n, p.spec.StartHeightUnits, int(p.spec.StartDirection))
	// 	fmt.Printf("f(%d) = %d\n", n, y)
	// 	editor.DrawBar2(y, p.barPlacement)
	// }
	n := p.nextNum + barCount
	for x := p.nextNum; x < n; x++ {

		if common.IsPrime(x) {
			p.barPlacement = p.barPlacement.Toggle()
		}

		y := countuneFn(x-p.spec.StartNum, p.spec.StartHeightUnits, int(p.spec.StartDirection))
		// fmt.Printf("f(%d) = %d\n", x, y)
		editor.DrawBar2(y, p.barPlacement)
	}

	p.nextNum += barCount

	return editor.image()
}
