package tasks

type ModelData struct {
	SizeA int
	SizeB int
	Model [][]float64 // Матрица коэффициентов (рецепты)
	A []float64   // Вектор ресурсов (склад)
	B []float64   // Вектор плана (что хотим испечь)
}

var CurrentData = ModelData{
	SizeA: 0,
	SizeB: 0,
	Model: [][]float64{},
	A: []float64{},
	B: []float64{},
}