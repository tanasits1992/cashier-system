package services

import "sort"

type SeriesService struct{}

func NewSeriesService() *SeriesService {
	return &SeriesService{}
}

func (s *SeriesService) Calculate(positions []int) map[int]int {
	sort.Slice(positions, func(i, j int) bool {
		return positions[i] < positions[j]
	})

	cur := 1
	seriesMap := make(map[int]int)
	inc1 := 0
	inc2 := 1

	for i := 0; i < positions[len(positions)-1]; i++ {
		seriesMap[i+1] = cur

		inc2 = inc2 + 1
		inc1 = inc1 + inc2
		cur = cur + inc1

	}

	result := map[int]int{}
	for _, val := range positions {
		result[val] = seriesMap[val]
	}

	return result
}
