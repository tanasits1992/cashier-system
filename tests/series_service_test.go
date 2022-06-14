package tests

import (
	"cashier/app/services"
	"strconv"
	"strings"
	"testing"
)

var mockSeriesService services.SeriesService

func mockTestSeriesServiceSetup(m *testing.T) {
	mockSeriesService = *services.NewSeriesService()
}

func mockTestSeriesServiceShutdown(m *testing.T) {
}

func TestCalculateSeries(t *testing.T) {
	cases := []struct {
		name   string
		input  []int
		result string
	}{
		{
			name:   "success 1",
			input:  []int{1, 2, 3},
			result: "1,3,8",
		},
		{
			name:   "success 2",
			input:  []int{7, 10, 12},
			result: "78,211,353",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockTestSeriesServiceSetup(t)
			defer mockTestSeriesServiceShutdown(t)

			series := mockSeriesService.Calculate(c.input)
			result := []string{}

			for _, num := range c.input {
				result = append(result, strconv.Itoa(series[num]))
			}

			Assert(t, "series", c.result, strings.Join(result, ","))

		})
	}
}
