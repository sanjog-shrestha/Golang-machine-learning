// Package mlmodel implements simple linear regression from scratch.
package mlmodel

import (
	"fmt"
	"math"
)

// ===== 1. DEFINE MODEL AND DATA =====

// Sample is one training example: a feature x and its true label y.
type Sample struct {
	X float64 // feature, e.g. matches played
	Y float64 // target,  e.g. goals scored
}

// LinearModel holds the learnable parameters of y = w*x + b.
type LinearModel struct {
	Weight float64
	Bias   float64
}

// ===== 2. HYPOTHESIS FUNCTION =====

// Predict applies the hypothesis h(x) = w*x + b to a single input.
func (m LinearModel) Predict(x float64) float64 {
	return m.Weight*x + m.Bias
}

// ===== 3. COST FUNCTION =====

// Cost computes Mean Squared Error (MSE) over the dataset:
//
//	J(w,b) = (1/2n) * Σ (h(xᵢ) - yᵢ)²
//
// The 1/2 factor makes the gradient cleaner.
func (m LinearModel) Cost(data []Sample) float64 {
	n := len(data)
	if n == 0 {
		return 0
	}
	var sum float64
	for _, s := range data {
		err := m.Predict(s.X) - s.Y
		sum += err * err
	}
	return sum / (2 * float64(n))
}

// ===== 4. GRADIENT DESCENT (one step) =====

// gradientStep nudges Weight and Bias one step down the cost surface.
// Partial derivatives of MSE:
//
//	∂J/∂w = (1/n) Σ (h(xᵢ) - yᵢ) * xᵢ
//	∂J/∂b = (1/n) Σ (h(xᵢ) - yᵢ)
func (m *LinearModel) gradientStep(data []Sample, lr float64) {
	n := float64(len(data))
	var dW, dB float64
	for _, s := range data {
		err := m.Predict(s.X) - s.Y
		dW += err
		dB += err
	}
	dW /= n
	dB /= n

	// update parameters in the opposite direction of the gradient
	m.Weight -= lr * dW
	m.Bias -= lr * dB
}

// ===== 5. TRAIN THE MODEL =====

// Train runs gradient descent for a number of epochs.
// It returns the cost history so you can watch it converge.
func (m *LinearModel) Train(data []Sample, lr float64, epochs int) []float64 {
	history := make([]float64, 0, epochs)
	for i := 0; i < epochs; i++ {
		m.gradientStep(data, lr)
		history = append(history, m.Cost(data))
	}
	return history
}

// ===== 6. EVALUATION =====

// RSquared (coefficient of determination) measures goodness of fit.
// 1.0 = perfect fit; 0 = no better than predicting the mean.
func (m *LinearModel) RSquared(data []Sample) float64 {
	n := len(data)
	if n == 0 {
		return 0
	}
	var meanY float64
	for _, s := range data {
		meanY += s.Y
	}
	meanY /= float64(n)

	var ssRes, ssTot float64 // residual and total sum of squares

	for _, s := range data {
		pred := m.Predict(s.X)
		ssRes += (s.Y - pred) * (s.Y - pred)
		ssTot += (s.Y - meanY) * (s.Y - meanY)
	}
	if ssTot == 0 {
		return 0
	}
	return 1 - ssRes/ssTot
}

func (m LinearModel) RMSE(data []Sample) float64 {
	n := len(data)
	if n == 0 {
		return 0
	}
	var sum float64
	for _, s := range data {
		err := m.Predict(s.X) - s.Y
		sum += err * err
	}
	return math.Sqrt(sum / float64(n))
}

func (m LinearModel) Summary(data []Sample) {
	fmt.Printf("Learned model: goals = %.3f * matches + %.3f\n", m.Weight, m.Bias)
	fmt.Printf("  R²   : %.3f\n", m.RSquared(data))
	fmt.Printf("  RMSE : %.3f\n", m.RMSE(data))
}

func Normalize(data []Sample) (scaled []Sample, min, span float64) {
	if len(data) == 0 {
		return nil, 0, 1
	}
	min, max := data[0].X, data[0].X
	for _, s := range data {
		if s.X < min {
			min = s.X
		}
		if s.X > max {
			max = s.X
		}
	}
	span = max - min
	if span == 0 {
		span = 1
	}
	scaled = make([]Sample, len(data))
	for i, s := range data {
		scaled[i] = Sample{X: (s.X - min) / span, Y: s.Y}
	}
	return scaled, min, span
}

// ===== MEAN SQUARED ERROR (evaluation metric) =====

// MSE is the standard Mean Squared Error: (1/n) Σ (h(xᵢ) - yᵢ)².
// Note: this differs from Cost(), which uses 1/2n for cleaner gradients.
// MSE is what you report; Cost is what you optimize.
func (m LinearModel) MSE(data []Sample) float64 {
	n := len(data)
	if n == 0 {
		return 0
	}
	var sum float64
	for _, s := range data {
		err := m.Predict(s.X) - s.Y
		sum += err * err
	}
	return sum / float64(n)
}

// ===== RESIDUALS (diagnostic) =====

// Residual pairs a feature value with its prediction error (y - ŷ).
type Residual struct {
	X        float64
	Residual float64
}

// Residuals computes y - prediction for each sample.
// Good fit -> residuals scattered randomly around zero.
// Pattern/curve in residuals -> a linear model is the wrong choice.
func (m LinearModel) Residuals(data []Sample) []Residual {
	out := make([]Residual, len(data))
	for i, s := range data {
		out[i] = Residual{X: s.X, Residual: s.Y - m.Predict(s.X)}
	}
	return out
}

// ===== L1 REGULARIZATION (Lasso) =====

// CostL1 is MSE-style cost plus the L1 penalty: (lambda/n) * |w|.
// The bias term is conventionally NOT regularized.
func (m LinearModel) CostL1(data []Sample, lambda float64) float64 {
	base := m.Cost(data) // existing 1/2n MSE cost
	n := float64(len(data))
	if n == 0 {
		return base
	}

	return base + (lambda/n)*math.Abs(m.Weight)
}

// sign returns -1, 0, or +1 — the derivative of |w|.
func sign(x float64) float64 {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}

// gradientStepL1 is gradient descent with an L1 penalty on the weight.
// The penalty's gradient is lambda * sign(w), applied to the weight only.
func (m *LinearModel) gradientStepL1(data []Sample, lr, lambda float64) {
	n := float64(len(data))
	var dW, dB float64
	for _, s := range data {
		err := m.Predict(s.X) - s.Y
		dW += err * s.X
		dB += err
	}
	dW = dW/n + (lambda/n)*sign(m.Weight) // L1 term added to weight gradient
	dB /= n                               // bias left unregularized

	m.Weight -= lr * dW
	m.Bias -= lr * dB
}

// TrainL1 runs regularized gradient descent and returns the cost history.
func (m *LinearModel) TrainL1(data []Sample, lr, lambda float64, epochs int) []float64 {
	history := make([]float64, 0, epochs)
	for i := 0; i < epochs; i++ {
		m.gradientStepL1(data, lr, lambda)
		history = append(history, m.CostL1(data, lambda))
	}
	return history
}
