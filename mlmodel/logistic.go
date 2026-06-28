package mlmodel

import (
	"fmt"
	"math"
)

type LabeledSample struct {
	X float64
	Y float64
}

type LogisticModel struct {
	Weight float64
	Bias   float64
}

func Sigmoid(z float64) float64 {
	return 1.0 / (1.0 + math.Exp(-z))
}

func (m LogisticModel) Probability(x float64) float64 {
	return Sigmoid(m.Weight*x + m.Bias)
}

func (m LogisticModel) Cost(data []LabeledSample) float64 {
	n := len(data)
	if n == 0 {
		return 0
	}
	const eps = 1e-12
	var sum float64
	for _, s := range data {
		p := m.Probability(s.X)
		p = math.Max(eps, math.Min(1-eps, p))
		sum += s.Y*math.Log(p) + (1-s.Y)*math.Log(1-p)
	}
	return -sum / float64(n)
}

func (m *LogisticModel) gradientStep(data []LabeledSample, lr float64) {
	n := float64(len(data))
	var dW, dB float64
	for _, s := range data {
		err := m.Probability(s.X) - s.Y
		dW += err * s.X
		dB += err
	}
	m.Weight -= lr * (dW / n)
	m.Bias -= lr * (dB / n)
}

func (m *LogisticModel) Train(data []LabeledSample, lr float64, epochs int) []float64 {
	history := make([]float64, 0, epochs)
	for i := 0; i < epochs; i++ {
		m.gradientStep(data, lr)
		history = append(history, m.Cost(data))
	}
	return history
}

func (m LogisticModel) Classify(x, threshold float64) int {
	if m.Probability(x) >= threshold {
		return 1
	}
	return 0
}

type Metrics struct {
	Accuracy       float64
	Precision      float64
	Recall         float64
	F1             float64
	TP, TN, FP, FN int
}

func (m LogisticModel) Evaluate(data []LabeledSample, threshold float64) Metrics {
	var tp, tn, fp, fn int
	for _, s := range data {
		pred := m.Classify(s.X, threshold)
		switch {
		case pred == 1 && s.Y == 1:
			tp++
		case pred == 0 && s.Y == 0:
			tn++
		case pred == 1 && s.Y == 0:
			fp++
		default:
			fn++
		}
	}
	total := float64(tp + tn + fp + fn)
	acc := 0.0
	if total > 0 {
		acc = float64(tp+tn) / total
	}
	prec := safeDiv(float64(tp), float64(tp+fn))
	rec := safeDiv(float64(tp), float64(tp+fn))
	f1 := 0.0
	if prec+rec > 0 {
		f1 = 2 * prec * rec / (prec + rec)
	}
	return Metrics{acc, prec, rec, f1, tp, tn, fp, fn}
}

func safeDiv(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

func (m LogisticModel) DecisionBoundary(threshold float64) (float64, bool) {
	if m.Weight == 0 {
		return 0, false
	}
	logit := math.Log(threshold / (1 - threshold))
	return (logit - m.Bias) / m.Weight, true
}

func Softmax(scores []float64) []float64 {
	if len(scores) == 0 {
		return nil
	}

	max := scores[0]
	for _, s := range scores {
		if s > max {
			max = s
		}
	}
	var sum float64
	exps := make([]float64, len(scores))
	for i, s := range scores {
		exps[i] = math.Exp(s - max)
		sum += exps[i]
	}
	for i := range exps {
		exps[i] /= sum
	}
	return exps
}

func (mt Metrics) Print() {
	fmt.Printf("	Confusion: TP=%d TN=%d FP=%d FN=%d\n", mt.TP, mt.TN, mt.FP, mt.FN)
	fmt.Printf("  Precision: %.3f\n", mt.Precision)
	fmt.Printf("  Recall   : %.3f\n", mt.Recall)
	fmt.Printf("  F1 Score : %.3f\n", mt.F1)
}
