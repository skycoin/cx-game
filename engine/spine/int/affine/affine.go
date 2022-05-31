package affine

// add is deprecated
func add(lhs, rhs []float64, dim int) []float64 {
	result := make([]float64, len(lhs))
	for i := 0; i < dim-1; i++ {
		for j := 0; j < dim; j++ {
			result[i*dim+j] = lhs[i*dim+j] + rhs[i*dim+j]
		}
	}
	return result
}

func mulSquare(lhs, rhs []float32, dim int) []float32 {
	result := make([]float32, len(lhs))
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			e := float32(0.0)
			for k := 0; k < dim; k++ {
				e += lhs[i*dim+k] * rhs[k*dim+j]
			}
			result[i*dim+j] = e
		}
	}
	return result
}
