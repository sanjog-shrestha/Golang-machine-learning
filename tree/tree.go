// Package tree implements decision trees and random forests from scratch
// for binary classification.
package tree

import (
	"math/rand"
	"sort"
)

// ===== 1. DATA =====

// Row is one labeled example with multiple features and a binary class.
type Row struct {
	Features []float64 // e.g. [matches, goals]
	Label    int       // 0 or 1
}

// ===== 2. IMPURITY: GINI =====

// gini measures how mixed the labels are in a group.
// 0 = pure (all one class); 0.5 = maximally mixed for two classes.
//
//	Gini = 1 - Σ pᵢ²
func gini(rows []Row) float64 {
	n := len(rows)
	if n == 0 {
		return 0
	}
	counts := map[int]int{}
	for _, r := range rows {
		counts[r.Label]++
	}
	impurity := 1.0
	for _, c := range counts {
		p := float64(c) / float64(n)
		impurity -= p * p
	}
	return impurity
}

// ===== 3. SPLITTING =====

// split partitions rows by whether feature[idx] <= threshold.
func split(rows []Row, idx int, threshold float64) (left, right []Row) {
	for _, r := range rows {
		if r.Features[idx] <= threshold {
			left = append(left, r)
		} else {
			right = append(right, r)
		}
	}
	return
}

// bestSplit searches features and thresholds for the split that most
// reduces weighted Gini impurity (the largest "information gain").
// featureSubset limits which features are considered — used by the forest.
func bestSplit(rows []Row, featureSubset []int) (bestIdx int, bestThr, bestGain float64) {
	current := gini(rows)
	n := float64(len(rows))
	bestIdx = -1

	for _, idx := range featureSubset {
		// candidate thresholds = sorted unique values of this feature
		vals := make([]float64, len(rows))
		for i, r := range rows {
			vals[i] = r.Features[idx]
		}
		sort.Float64s(vals)

		for _, thr := range vals {
			left, right := split(rows, idx, thr)
			if len(left) == 0 || len(right) == 0 {
				continue
			}
			// weighted impurity of the two children
			wImpurity := (float64(len(left))/n)*gini(left) + (float64(len(right))/n)*gini(right)
			gain := current - wImpurity
			if gain > bestGain {
				bestGain, bestIdx, bestThr = gain, idx, thr
			}
		}

	}
	return bestIdx, bestThr, bestGain
}

// ===== 4. THE TREE =====

// Node is one node in the tree: either a leaf (Prediction set) or an
// internal split (FeatureIdx/Threshold with Left/Right children).
type Node struct {
	FeatureIdx int
	Threshold  float64
	Left       *Node
	Right      *Node
	Leaf       bool
	Prediction int
}

// DecisionTree wraps a root node plus stopping hyperparameters.
type DecisionTree struct {
	Root        *Node
	MaxDepth    int
	MinLeafSize int
	// featuresPerSplit: if > 0, only a random subset of features is
	// considered at each split (random forests use this; a lone tree sets 0).
	featuresPerSplit int
}

// majority returns the most common label in rows (the leaf's prediction).
func majority(rows []Row) int {
	counts := map[int]int{}
	for _, r := range rows {
		counts[r.Label]++
	}
	best, bestCount := 0, -1
	for label, c := range counts {
		if c > bestCount {
			best, bestCount = label, c
		}
	}
	return best
}

// build recursively grows the tree until a stopping condition is hit.
func (t *DecisionTree) build(rows []Row, depth int) *Node {
	// stopping conditions: pure node, too deep, or too few samples
	if len(rows) == 0 {
		return &Node{Leaf: true, Prediction: 0}
	}
	if gini(rows) == 0 || depth >= t.MaxDepth || len(rows) <= t.MinLeafSize {
		return &Node{Leaf: true, Prediction: majority(rows)}
	}

	// choose which features to consider at this split
	nFeatures := len(rows[0].Features)
	subset := featureSubset(nFeatures, t.featuresPerSplit)

	idx, thr, gain := bestSplit(rows, subset)
	if idx == -1 || gain <= 0 {
		return &Node{Leaf: true, Prediction: majority(rows)}
	}

	left, right := split(rows, idx, thr)
	return &Node{
		FeatureIdx: idx,
		Threshold:  thr,
		Left:       t.build(left, depth+1),
		Right:      t.build(right, depth+1),
	}
}

// Fit trains the tree on the given rows.'
func (t *DecisionTree) Fit(rows []Row) {
	if t.MaxDepth == 0 {
		t.MaxDepth = 10
	}
	if t.MinLeafSize == 0 {
		t.MinLeafSize = 1
	}
	t.Root = t.build(rows, 0)
}

// predictNode walks the tree to a leaf for one row.
func predictNode(n *Node, features []float64) int {
	if n.Leaf {
		return n.Prediction
	}
	if features[n.FeatureIdx] <= n.Threshold {
		return predictNode(n.Left, features)
	}
	return predictNode(n.Right, features)
}

// Predict classifies a single feature vector.
func (t *DecisionTree) Predict(features []float64) int {
	return predictNode(t.Root, features)
}

// featureSubset returns either all feature indices (k<=0) or k random ones.
func featureSubset(nFeatures, k int) []int {
	all := make([]int, nFeatures)
	for i := range all {
		all[i] = i
	}
	if k <= 0 || k >= nFeatures {
		return all
	}
	rand.Shuffle(len(all), func(i, j int) { all[i], all[j] = all[j], all[i] })
	return all[:k]
}
