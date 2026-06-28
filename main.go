package main

import (
	"encoding/json"
	"fmt"
	"os"

	"go-sports/athlete"
	"go-sports/mlmodel"
	"go-sports/preprocess"
	"go-sports/stats" // NEW
	"go-sports/viz"   // NEW
)

type Roster struct {
	Footballers []athlete.Footballer `json:"footballers"`
	Cricketers  []athlete.Cricketer  `json:"cricketers"`
}

func main() {
	data, err := os.ReadFile("players.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var roster Roster
	if err := json.Unmarshal(data, &roster); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// --- Cleansing (from previous step) ---
	roster.Footballers = preprocess.CleanFootballers(roster.Footballers)
	roster.Cricketers = preprocess.CleanCricketers(roster.Cricketers)

	// ============ EXPLORATORY DATA ANALYSIS ============
	fmt.Println("===== EDA: Football Goals =====")

	// extract the numeric field we want to analyze
	goals := make([]float64, 0, len(roster.Footballers))
	for _, f := range roster.Footballers {
		goals = append(goals, float64(f.Goals))
	}
	goalSummary := stats.Describe(goals)
	goalSummary.Print("Goals")

	// categorical EDA: how many players per team
	teams := make([]string, 0)
	for _, f := range roster.Footballers {
		teams = append(teams, f.Team)
	}
	freq := stats.Frequency(teams)
	fmt.Println("\n--- Players per Team ---")
	for team, n := range freq {
		fmt.Printf("  %s: %d\n", team, n)
	}

	// ============ DATA VISUALIZATION ============
	// Build a bar chart of goals per footballer.
	bars := make([]viz.Bar, 0, len(roster.Footballers))
	for _, f := range roster.Footballers {
		bars = append(bars, viz.Bar{Label: f.Name, Value: float64(f.Goals)})
	}

	html := viz.BarChartHTML("Goals by Footballer", bars)
	if err := viz.WriteHTML("goals_chart.html", html); err != nil {
		fmt.Println("Error writing chart:", err)
		return
	}
	fmt.Println("\nChart written to goals_chart.html (open it in a browser)")

	// ============ LINEAR REGRESSION ============
	fmt.Println("\n==== Linear Regression: predict goals from matches ====")

	// 1. Build the dataset. In a real project you'd pull these from your
	//    cleaned roster; here we use illustrative (matches, goals) pairs.
	raw := []mlmodel.Sample{
		{X: 10, Y: 4},
		{X: 20, Y: 9},
		{X: 30, Y: 14},
		{X: 40, Y: 18},
		{X: 50, Y: 23},
	}

	// 2. Normalize features so gradient descent converges reliably.
	scaled, min, span := mlmodel.Normalize(raw)

	// 3. Train.
	model := mlmodel.LinearModel{} // starts at w=0, b=0
	history := model.Train(scaled, 0.1, 1000)

	// 4. Report training progress and final fit.
	fmt.Printf("Cost: start=%.4f  end=%.4f\n", history[0], history[len(history)-1])
	model.Summary(scaled)

	// 5. Make a prediction for a new player with 35 matches.
	newMatches := 35.0
	pred := model.Predict((newMatches - min) / span)
	fmt.Printf("\nPrediction: a player with %.0f matches scores ~%.1f goals\n", newMatches, pred)

	// ============ FINE-TUNING: evaluation & regularization ============
	fmt.Println("\n==== Fine-Tuning ====")

	fmt.Printf("MSE (no regularization): %.4f\n", model.MSE(scaled))

	lassoModel := mlmodel.LinearModel{}
	lambda := 0.1
	lassoModel.TrainL1(scaled, 0.1, lambda, 1000)
	fmt.Printf("L1 (lambda=%.2f): weight:%.3f bias=%.3f MSE=%.4f\n", lambda, lassoModel.Weight, lassoModel.Bias, lassoModel.MSE(scaled))

	res := model.Residuals(scaled)
	pts := make([]viz.Point, len(res))
	for i, r := range res {
		pts[i] = viz.Point{X: r.X, Y: r.Residual}
	}
	resHTML := viz.ResidualPlotHTML("Residuals: goals model", pts)
	if err := viz.WriteHTML("residuals.html", resHTML); err != nil {
		fmt.Println("Error writing residual plot", err)
	} else {
		fmt.Println("Residual plot written to residuals.html")
	}

	// ============ LOGISTIC REGRESSION: top-scorer classification ============
	fmt.Println("\n===== Logistic Regression: is this player a top scorer? =====")

	clsRaw := []mlmodel.LabeledSample{
		{X: 0.1, Y: 0},
		{X: 0.2, Y: 0},
		{X: 0.35, Y: 0},
		{X: 0.5, Y: 1},
		{X: 0.7, Y: 1},
		{X: 0.9, Y: 1},
	}

	clf := mlmodel.LogisticModel{}
	clsHistory := clf.Train(clsRaw, 0.5, 5000)
	fmt.Printf("Log-loss: start=%.4f end=%.4f\n", clsHistory[0], clsHistory[len(clsHistory)-1])

	metrics := clf.Evaluate(clsRaw, 0.5)
	metrics.Print()

	if x, ok := clf.DecisionBoundary(0.5); ok {
		fmt.Printf("	Decision boundary at x = %.3f\n", x)
	}

	p := 0.6
	fmt.Printf("\nPlayer at x=%2f: P(top scorer)=%.3f -> class %d\n", p, clf.Probability(p), clf.Classify(p, 0.5))

	scores := []float64{2.0, 1.0, 0.1}
	probs := mlmodel.Softmax(scores)
	fmt.Printf("\nMulti-class Softmax %v -> %.3f, %.3f, %.3f\n", scores, probs[0], probs[1], probs[2])
}
