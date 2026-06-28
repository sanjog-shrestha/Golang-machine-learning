package tree

import (
	"math"
	"math/rand"
)

// ===== 5. RANDOM FOREST =====

// RandomForest is an ensemble of decision trees trained on bootstrap
// samples, each split considering a random subset of features.
// Predictions are made by majority vote across all trees.
type RandomForest struct {
	Trees       []*DecisionTree
	NTrees      int
	MaxDepth    int
	MinLeafSize int
}

// bootstrap samples len(rows) rows WITH replacement (bagging).
// On average ~63% of rows appear; some repeat, some are left out.
func bootstrap(rows []Row) []Row {
	n := len(rows)
	sample := make([]Row, n)
	for i := 0; i < n; i++ {
		sample[i] = rows[rand.Intn(n)]
	}
	return sample
}

// Fit trains all trees. Each tree gets its own bootstrap sample and
// considers sqrt(nFeatures) features per split — the standard default.
func (rf *RandomForest) Fit(rows []Row) {
	if rf.NTrees == 0 {
		rf.NTrees = 10
	}
	if rf.MaxDepth == 0 {
		rf.MaxDepth = 8
	}
	nFeatures := len(rows[0].Features)
	featuresPerSplit := int(math.Sqrt(float64(nFeatures)))
	if featuresPerSplit < 1 {
		featuresPerSplit = 1
	}
	rf.Trees = make([]*DecisionTree, rf.NTrees)
	for i := 0; i < rf.NTrees; i++ {
		t := &DecisionTree{
			MaxDepth:         rf.MaxDepth,
			MinLeafSize:      rf.MinLeafSize,
			featuresPerSplit: featuresPerSplit,
		}
		t.Fit(bootstrap(rows)) // each tree sees a different resampled dataset
		rf.Trees[i] = t
	}
}

// Predict takes a majority vote across all trees.
func (rf *RandomForest) Predict(features []float64) int {
	votes := map[int]int{}
	for _, t := range rf.Trees {
		votes[t.Predict(features)]++
	}
	best, bestCount := 0, -1
	for label, c := range votes {
		if c > bestCount {
			best, bestCount = label, c
		}
	}
	return best
}

// Accuracy evaluates the forest on labeled data.
func (rf *RandomForest) Accuracy(rows []Row) float64 {
	if len(rows) == 0 {
		return 0
	}
	correct := 0
	for _, r := range rows {
		if rf.Predict(r.Features) == r.Label {
			correct++
		}
	}
	return float64(correct) / float64(len(rows))
}
