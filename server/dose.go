package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/kennedyjustin/BolusGPT/bolus"
)

type DoseInput struct {
	TotalGramsOfCarbs   float32 `json:"total_grams_of_carbs"`
	GramsOfFiber        float32 `json:"grams_of_fiber"`
	GramsOfSugarAlcohol float32 `json:"grams_of_sugar_alcohol"`
	GramsOfProtein      float32 `json:"grams_of_protein"`

	MinutesOfExercise float32                 `json:"minutes_of_exercise"`
	ExerciseIntensity bolus.ExerciseIntensity `json:"exercise_intensity"`
}

func (s *Server) DoseHandler(response http.ResponseWriter, request *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var dose bolus.Dose
	s.db.Read(func(data *Me) {
		if data == nil {
			http.Error(response, "please onboard", http.StatusNotFound)
			return
		}

		if data.InsulinToCarbRatio == 0 || data.InsulinSensitivityFactor == 0 || data.TargetBloodGlucoseLevelInMgDl == 0 {
			http.Error(response, "'insulin_to_carb_ratio', 'insulin_sensitivity_factor', and 'target_blood_glucose_level_in_mg_dl' required", http.StatusNotFound)
			return
		}

		decoder := json.NewDecoder(request.Body)
		input := DoseInput{}
		err := decoder.Decode(&input)
		if err != nil {
			log.Println(err)
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		currentBloodGlucoseReading, err := s.dexcomClient.GetCurrentBloodGlucoseReading()
		if err != nil {
			log.Println(err)
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		dose = bolus.GetDose(bolus.DoseInput{
			FoodInput: bolus.FoodInput{
				TotalGramsOfCarbs:                input.TotalGramsOfCarbs,
				GramsOfFiber:                     input.GramsOfFiber,
				FiberMultiplier:                  data.FiberMultiplier,
				GramsOfSugarAlcohol:              input.GramsOfSugarAlcohol,
				SugarAlcoholMultiplier:           data.SugarAlcoholMultiplier,
				GramsOfProtein:                   input.GramsOfProtein,
				ProteinMultiplier:                data.ProteinMultiplier,
				CarbThresholdToCountProteinUnder: data.CarbThresholdToCountProteinUnder,
				InsulinToCarbRatio:               data.InsulinToCarbRatio,
			},
			CorrectionInput: bolus.CorrectionInput{
				CurrentBloodGlucoseLevelInMgDl:  float32(currentBloodGlucoseReading.Value),
				BloodGlucoseTrendInMgDlIn15Mins: float32(currentBloodGlucoseReading.Get15MinDeltaFromTrend()),
				TargetBloodGlucoseLevelInMgDl:   data.TargetBloodGlucoseLevelInMgDl,
				InsulinSensitivityFactor:        data.InsulinSensitivityFactor,
			},
			InsulinOnBoardInput: bolus.InsulinOnBoardInput{
				LastBolusTime:           data.LastBolusTime,
				LastBolusUnitsOfInsulin: data.LastBolusUnitsOfInsulin,
			},
			ExerciseInput: bolus.ExerciseInput{
				MinutesOfExercise: input.MinutesOfExercise,
				ExerciseIntensity: input.ExerciseIntensity,
			},
		})
	})

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(dose)
}

type ConfirmInput struct {
	UnitsOfInsulin float32 `json:"units_of_insulin"`
}

func (s *Server) ConfirmHandler(response http.ResponseWriter, request *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	decoder := json.NewDecoder(request.Body)
	input := ConfirmInput{}
	err := decoder.Decode(&input)
	if err != nil {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.db.Write(func(me *Me) error {
		me.LastBolusUnitsOfInsulin = input.UnitsOfInsulin
		me.LastBolusTime = time.Now()

		return nil
	})
	if err != nil {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
}
