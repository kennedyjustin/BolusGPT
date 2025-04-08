package bolus

import (
	"testing"
	"time"
)

func TestDoseFoodFactor(t *testing.T) {
	dose := GetDose(DoseInput{
		FoodInput: FoodInput{
			TotalGramsOfCarbs:                20,
			GramsOfFiber:                     8,
			FiberMultiplier:                  0.5,
			GramsOfSugarAlcohol:              10,
			SugarAlcoholMultiplier:           0.5,
			GramsOfProtein:                   18,
			ProteinMultiplier:                0.5,
			CarbThresholdToCountProteinUnder: 21,
			InsulinToCarbRatio:               SimpleTimeSensitiveFactor(5),
		},
		CorrectionInput: CorrectionInput{
			CurrentBloodGlucoseLevelInMgDl: 100,
			TargetBloodGlucoseLevelInMgDl:  100,
			InsulinSensitivityFactor:       SimpleTimeSensitiveFactor(30),
		},
	})
	if dose.UnitsOfInsulin != 4 {
		t.Errorf("expected 4, got %f", dose.UnitsOfInsulin)
	}
	if dose.GramsOfCarbs != 0 {
		t.Errorf("expected 0, got %f", dose.GramsOfCarbs)
	}
	if dose.Breakdown.FoodFactor != 4 {
		t.Errorf("expected 4, got %f", dose.Breakdown.FoodFactor)
	}
	if dose.Breakdown.CorrectionFactor != 0 {
		t.Errorf("expected 0, got %f", dose.Breakdown.CorrectionFactor)
	}
	if dose.Breakdown.InsulinOnBoardFactor != 0 {
		t.Errorf("expected 0, got %f", dose.Breakdown.InsulinOnBoardFactor)
	}
	if dose.Breakdown.ExerciseMultiplier != 1 {
		t.Errorf("expected 1, got %f", dose.Breakdown.ExerciseMultiplier)
	}
}

func TestDoseCorrectionFactor(t *testing.T) {
	dose := GetDose(DoseInput{
		FoodInput: FoodInput{
			InsulinToCarbRatio: SimpleTimeSensitiveFactor(5),
		},
		CorrectionInput: CorrectionInput{
			CurrentBloodGlucoseLevelInMgDl:  120,
			BloodGlucoseTrendInMgDlIn15Mins: 15,
			TargetBloodGlucoseLevelInMgDl:   90,
			InsulinSensitivityFactor:        SimpleTimeSensitiveFactor(30),
		},
	})
	if dose.UnitsOfInsulin != 1.5 {
		t.Errorf("expected 1.5, got %f", dose.UnitsOfInsulin)
	}
	if dose.GramsOfCarbs != 0 {
		t.Errorf("expected 0, got %f", dose.GramsOfCarbs)
	}
	if dose.Breakdown.FoodFactor != 0 {
		t.Errorf("expected 4, got %f", dose.Breakdown.FoodFactor)
	}
	if dose.Breakdown.CorrectionFactor != 1.5 {
		t.Errorf("expected 1.5, got %f", dose.Breakdown.CorrectionFactor)
	}
	if dose.Breakdown.InsulinOnBoardFactor != 0 {
		t.Errorf("expected 0, got %f", dose.Breakdown.InsulinOnBoardFactor)
	}
	if dose.Breakdown.ExerciseMultiplier != 1 {
		t.Errorf("expected 1, got %f", dose.Breakdown.ExerciseMultiplier)
	}
}

func TestDoseInsulinOnBoardFactor(t *testing.T) {
	dose := GetDose(DoseInput{
		FoodInput: FoodInput{
			InsulinToCarbRatio: SimpleTimeSensitiveFactor(5),
		},
		CorrectionInput: CorrectionInput{
			CurrentBloodGlucoseLevelInMgDl: 190,
			TargetBloodGlucoseLevelInMgDl:  100,
			InsulinSensitivityFactor:       SimpleTimeSensitiveFactor(30),
		},
		InsulinOnBoardInput: InsulinOnBoardInput{
			LastBolusTime:           time.Now().Add(-100 * time.Minute),
			LastBolusUnitsOfInsulin: 2,
		},
	})
	if dose.UnitsOfInsulin != 2 {
		t.Errorf("expected 2, got %f", dose.UnitsOfInsulin)
	}
	if dose.GramsOfCarbs != 0 {
		t.Errorf("expected 0, got %f", dose.GramsOfCarbs)
	}
	if dose.Breakdown.FoodFactor != 0 {
		t.Errorf("expected 0, got %f", dose.Breakdown.FoodFactor)
	}
	if dose.Breakdown.CorrectionFactor != 3 {
		t.Errorf("expected 3, got %f", dose.Breakdown.CorrectionFactor)
	}
	if dose.Breakdown.InsulinOnBoardFactor != -1 {
		t.Errorf("expected -1, got %f", dose.Breakdown.InsulinOnBoardFactor)
	}
	if dose.Breakdown.ExerciseMultiplier != 1 {
		t.Errorf("expected 1, got %f", dose.Breakdown.ExerciseMultiplier)
	}
}

func TestExerciseMultiplier(t *testing.T) {
	dose := GetDose(DoseInput{
		FoodInput: FoodInput{
			InsulinToCarbRatio: SimpleTimeSensitiveFactor(5),
		},
		CorrectionInput: CorrectionInput{
			CurrentBloodGlucoseLevelInMgDl: 160,
			TargetBloodGlucoseLevelInMgDl:  100,
			InsulinSensitivityFactor:       SimpleTimeSensitiveFactor(30),
		},
		ExerciseInput: ExerciseInput{
			MinutesOfExercise: 45,
			ExerciseIntensity: High,
		},
	})
	if dose.UnitsOfInsulin != 1 {
		t.Errorf("expected 1, got %f", dose.UnitsOfInsulin)
	}
	if dose.GramsOfCarbs != 0.0 {
		t.Errorf("expected 0, got %f", dose.GramsOfCarbs)
	}
	if dose.Breakdown.FoodFactor != 0 {
		t.Errorf("expected 0, got %f", dose.Breakdown.FoodFactor)
	}
	if dose.Breakdown.CorrectionFactor != 2 {
		t.Errorf("expected 2, got %f", dose.Breakdown.CorrectionFactor)
	}
	if dose.Breakdown.InsulinOnBoardFactor != 0 {
		t.Errorf("expected 0, got %f", dose.Breakdown.InsulinOnBoardFactor)
	}
	if dose.Breakdown.ExerciseMultiplier != 0.5 {
		t.Errorf("expected 0.5, got %f", dose.Breakdown.ExerciseMultiplier)
	}
}

func TestGramsOfCarbsDueToLow(t *testing.T) {
	dose := GetDose(DoseInput{
		FoodInput: FoodInput{
			InsulinToCarbRatio: SimpleTimeSensitiveFactor(5),
		},
		CorrectionInput: CorrectionInput{
			CurrentBloodGlucoseLevelInMgDl:  70,
			BloodGlucoseTrendInMgDlIn15Mins: -15,
			TargetBloodGlucoseLevelInMgDl:   100,
			InsulinSensitivityFactor:        SimpleTimeSensitiveFactor(15),
		},
	})
	if dose.UnitsOfInsulin != -3 {
		t.Errorf("expected -3, got %f", dose.UnitsOfInsulin)
	}
	if dose.GramsOfCarbs != 15 {
		t.Errorf("expected 15, got %f", dose.GramsOfCarbs)
	}
	if dose.Breakdown.FoodFactor != 0 {
		t.Errorf("expected 0, got %f", dose.Breakdown.FoodFactor)
	}
	if dose.Breakdown.CorrectionFactor != -3 {
		t.Errorf("expected -3, got %f", dose.Breakdown.CorrectionFactor)
	}
	if dose.Breakdown.InsulinOnBoardFactor != 0 {
		t.Errorf("expected 0, got %f", dose.Breakdown.InsulinOnBoardFactor)
	}
	if dose.Breakdown.ExerciseMultiplier != 1 {
		t.Errorf("expected 1, got %f", dose.Breakdown.ExerciseMultiplier)
	}
}
