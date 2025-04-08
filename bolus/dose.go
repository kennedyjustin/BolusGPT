package bolus

import (
	"time"
)

type DoseInput struct {
	FoodInput
	CorrectionInput
	InsulinOnBoardInput
	ExerciseInput
}

type FoodInput struct {
	// Total (not net) Carbohydrates
	TotalGramsOfCarbs float32
	// Grams of Fiber
	GramsOfFiber float32
	// 0 to subtract none of the Fiber, 1 to subtract all of the Fiber
	FiberMultiplier float32
	// Grams of Sugar Alcohol
	GramsOfSugarAlcohol float32
	// 0 to subtract none of the Sugar Alcohol, 1 to subtract all of the Sugar Alcohol
	SugarAlcoholMultiplier float32
	// Grams of Protein
	GramsOfProtein float32
	// 0 to count none of the Protein, 1 to count all of the Protein
	ProteinMultiplier float32
	// For Carbs under the threshold, Protein will be counted
	CarbThresholdToCountProteinUnder float32
	// Insulin to Carb Ratio at a given time of day
	InsulinToCarbRatio TimeSensitiveFactor
}

type CorrectionInput struct {
	// Current Blood Sugar
	CurrentBloodGlucoseLevelInMgDl float32
	// Blood Sugar Trend (Delta)
	BloodGlucoseTrendInMgDlIn15Mins float32
	// Target Blood Sugar
	TargetBloodGlucoseLevelInMgDl float32
	// Insulin Sensitivity Factor at a given time of day
	InsulinSensitivityFactor TimeSensitiveFactor // Changes over time
}

type InsulinOnBoardInput struct {
	// Time of the last Bolus
	LastBolusTime time.Time
	// Units of Insulin for the last Bolus
	LastBolusUnitsOfInsulin float32
}

type ExerciseInput struct {
	// Minutes of Exercise post-Bolus
	MinutesOfExercise float32
	// Intensity of Exercise post-Bolus (Low, Medium, or High)
	ExerciseIntensity ExerciseIntensity
}

type TimeSensitiveFactor interface {
	GetAtTime(time.Time) float32
}

type SimpleTimeSensitiveFactor float32

func (s SimpleTimeSensitiveFactor) GetAtTime(time.Time) float32 {
	return float32(s)
}

type ExerciseIntensity int

const (
	Low ExerciseIntensity = iota
	Medium
	High
)

var InsulinOnBoardMultiplierList = []float32{
	1,    // 0 hours
	0.9,  // 0.5 hours
	0.7,  // 1 hour
	0.5,  // 1.5 hours
	0.35, // 2 hours
	0.2,  // 2.5 hours
	0.1,  // 3 hours
	0.05, // 3.5 hours
	0,    // 4 hours
}

var ExerciseMultiplierMap = []map[ExerciseIntensity]float32{
	{ // 0-30 minutes
		Low:    0.9,
		Medium: 0.75,
		High:   0.67,
	},
	{ // 30-60 minutes
		Low:    0.8,
		Medium: 0.67,
		High:   0.5,
	},
	{ // over 60 minutes
		Low:    0.7,
		Medium: 0.5,
		High:   0.33,
	},
}

type Dose struct {
	// Units of Insulin for the Bolus dose
	UnitsOfInsulin float32
	// If UnitsOfInsulin is negative, the grams of Carbohydrates to consume to get back to the Target Blood Glucose Level
	GramsOfCarbs float32
	// A breakdown of the major factors contributing to the Bolus dose
	Breakdown struct {
		FoodFactor           float32
		CorrectionFactor     float32
		InsulinOnBoardFactor float32
		ExerciseMultiplier   float32
	}
}

func GetDose(input DoseInput) Dose {
	dose := Dose{}

	// Validate Required Params
	if input.FoodInput.InsulinToCarbRatio == nil || input.CorrectionInput.TargetBloodGlucoseLevelInMgDl <= 55 || input.CorrectionInput.CurrentBloodGlucoseLevelInMgDl <= 0 || input.CorrectionInput.InsulinSensitivityFactor == nil {
		panic("valid inputs required for 'insulin_to_carb_ratio', 'target_blood_glucose_level_in_mg_dl', 'current_blood_glucose_level_in_mg_dl', and 'insulin_sensitivity_factor'")
	}

	// Calculate Food Factor
	grams := input.FoodInput.TotalGramsOfCarbs
	grams -= input.FoodInput.GramsOfFiber * input.FoodInput.FiberMultiplier
	grams -= input.FoodInput.GramsOfSugarAlcohol * input.FoodInput.SugarAlcoholMultiplier
	if grams < input.FoodInput.CarbThresholdToCountProteinUnder {
		grams += input.FoodInput.GramsOfProtein * input.FoodInput.ProteinMultiplier
	}
	insulinToCarbRatio := input.FoodInput.InsulinToCarbRatio.GetAtTime(time.Now())
	if insulinToCarbRatio > 0 {
		dose.Breakdown.FoodFactor = grams / input.FoodInput.InsulinToCarbRatio.GetAtTime(time.Now())
	}

	// Calculate Correction Factor
	bloodSugarIn15Mins := input.CorrectionInput.CurrentBloodGlucoseLevelInMgDl + input.CorrectionInput.BloodGlucoseTrendInMgDlIn15Mins
	correction := bloodSugarIn15Mins - input.CorrectionInput.TargetBloodGlucoseLevelInMgDl
	insulinSensitivityFactor := input.CorrectionInput.InsulinSensitivityFactor.GetAtTime(time.Now())
	if insulinSensitivityFactor > 0 {
		dose.Breakdown.CorrectionFactor = correction / input.CorrectionInput.InsulinSensitivityFactor.GetAtTime(time.Now())
	}

	// Calculate Insulin On Board
	incrementSinceLastBolus := int(time.Since(input.InsulinOnBoardInput.LastBolusTime).Minutes() / 30)
	if incrementSinceLastBolus < len(InsulinOnBoardMultiplierList) {
		dose.Breakdown.InsulinOnBoardFactor = input.InsulinOnBoardInput.LastBolusUnitsOfInsulin * -InsulinOnBoardMultiplierList[incrementSinceLastBolus]
	}

	// Calculate Exercise Multiplier
	if input.ExerciseInput.MinutesOfExercise > 0 {
		exerciseIncrement := int(input.ExerciseInput.MinutesOfExercise / 30)
		if exerciseIncrement < len(ExerciseMultiplierMap) {
			dose.Breakdown.ExerciseMultiplier = ExerciseMultiplierMap[exerciseIncrement][input.ExerciseInput.ExerciseIntensity]
		}
	}
	if dose.Breakdown.ExerciseMultiplier == 0 {
		dose.Breakdown.ExerciseMultiplier = 1
	}

	// Calculate Total. If negative, calculate the grams of carbs required to bring back to target.
	dose.UnitsOfInsulin = (dose.Breakdown.FoodFactor + dose.Breakdown.CorrectionFactor + dose.Breakdown.InsulinOnBoardFactor) * dose.Breakdown.ExerciseMultiplier
	if dose.UnitsOfInsulin < 0 && insulinToCarbRatio > 0 {
		dose.GramsOfCarbs = -dose.UnitsOfInsulin * insulinToCarbRatio
	}

	return dose
}
