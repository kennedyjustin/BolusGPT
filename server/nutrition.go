package server

import "net/http"

type NutritionSearch struct {
	Search string
}

type NutritionInfo struct {
	ServingSize         string
	TotalGramsOfCarbs   int
	GramsOfFiber        int
	GramsOfSugarAlcohol int
	GramsOfProtein      int
	GramsOfFat          int
}

func (s *Server) NutritionHandler(response http.ResponseWriter, request *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// TODO: Find Nutrition API
}
