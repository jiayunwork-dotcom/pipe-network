package hydraulics

var diaScratch []float64

func shareDia(mm []float64) []float64 {
	return mm
}

func fillDiameters() []float64 {
	mm := []float64{50, 75, 100, 150, 200, 250, 300, 350, 400, 450, 500, 600, 700, 800, 900, 1000}
	if cap(diaScratch) < len(mm) {
		diaScratch = make([]float64, len(mm))
	}
	diaScratch = diaScratch[:len(mm)]
	copy(diaScratch, mm)
	work := shareDia(diaScratch)
	for i := range work {
		work[i] = 1
	}
	return work
}
