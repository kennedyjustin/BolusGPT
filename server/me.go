package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/kennedyjustin/BolusGPT/bolus"
)

type Me struct {
	FiberMultiplier                  float32                         `json:"fiber_multiplier"`
	SugarAlcoholMultiplier           float32                         `json:"sugar_alcohol_multiplier"`
	ProteinMultiplier                float32                         `json:"protein_multiplier"`
	CarbThresholdToCountProteinUnder float32                         `json:"carb_threshold_to_count_protein_under"`
	InsulinToCarbRatio               bolus.SimpleTimeSensitiveFactor `json:"insulin_to_carb_ratio"`

	TargetBloodGlucoseLevelInMgDl float32                         `json:"target_blood_glucose_level_in_mg_dl"`
	InsulinSensitivityFactor      bolus.SimpleTimeSensitiveFactor `json:"insulin_sensitivity_factor"`

	LastBolusTime           time.Time `json:"last_bolus_time"`
	LastBolusUnitsOfInsulin float32   `json:"last_bolus_units_of_insulin"`
}

func (s *Server) MeHandlerGet(response http.ResponseWriter, request *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.db.Read(func(me *Me) {
		if me == nil {
			http.Error(response, "please onboard", http.StatusNotFound)
			return
		}

		response.Header().Set("Content-Type", "application/json")
		json.NewEncoder(response).Encode(me)
	})
}

type MeInput struct {
	FiberMultiplier                  *float32                         `json:"fiber_multiplier"`
	SugarAlcoholMultiplier           *float32                         `json:"sugar_alcohol_multiplier"`
	ProteinMultiplier                *float32                         `json:"protein_multiplier"`
	CarbThresholdToCountProteinUnder *float32                         `json:"carb_threshold_to_count_protein_under"`
	InsulinToCarbRatio               *bolus.SimpleTimeSensitiveFactor `json:"insulin_to_carb_ratio"`

	TargetBloodGlucoseLevelInMgDl *float32                         `json:"target_blood_glucose_level_in_mg_dl"`
	InsulinSensitivityFactor      *bolus.SimpleTimeSensitiveFactor `json:"insulin_sensitivity_factor"`

	LastBolusTime           *time.Time `json:"last_bolus_time"`
	LastBolusUnitsOfInsulin *float32   `json:"last_bolus_units_of_insulin"`
}

func (s *Server) MeHandlerPatch(response http.ResponseWriter, request *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	decoder := json.NewDecoder(request.Body)
	input := MeInput{}
	err := decoder.Decode(&input)
	if err != nil {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.db.Write(func(me *Me) error {
		if input.FiberMultiplier != nil {
			me.FiberMultiplier = *input.FiberMultiplier
		}
		if input.SugarAlcoholMultiplier != nil {
			me.SugarAlcoholMultiplier = *input.SugarAlcoholMultiplier
		}
		if input.ProteinMultiplier != nil {
			me.ProteinMultiplier = *input.ProteinMultiplier
		}
		if input.CarbThresholdToCountProteinUnder != nil {
			me.CarbThresholdToCountProteinUnder = *input.CarbThresholdToCountProteinUnder
		}
		if input.InsulinToCarbRatio != nil {
			me.InsulinToCarbRatio = *input.InsulinToCarbRatio
		}

		if input.TargetBloodGlucoseLevelInMgDl != nil {
			me.TargetBloodGlucoseLevelInMgDl = *input.TargetBloodGlucoseLevelInMgDl
		}
		if input.InsulinSensitivityFactor != nil {
			me.InsulinSensitivityFactor = *input.InsulinSensitivityFactor
		}

		if input.LastBolusTime != nil {
			me.LastBolusTime = *input.LastBolusTime
		}
		if input.LastBolusUnitsOfInsulin != nil {
			me.LastBolusUnitsOfInsulin = *input.LastBolusUnitsOfInsulin
		}

		response.Header().Set("Content-Type", "application/json")
		json.NewEncoder(response).Encode(me)

		return nil
	})
	if err != nil {
		log.Println(err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
}
