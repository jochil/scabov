package classifier

import (
	"math"
)

/*
	TODO document formula
*/
func QCorrelationCoefficient(rawMatrix map[string]map[string]float64) map[string]map[string]float64 {

	result := map[string]map[string]float64{}

	for k_name, k_props := range rawMatrix {
		result[k_name] = map[string]float64{}

		k_props_avg := calcPropsAvg(k_props)

		for l_name, l_props := range rawMatrix {

			l_props_avg := calcPropsAvg(l_props)

			var numerator float64 = 0
			var sum_k float64 = 0
			var sum_l float64 = 0

			for key, _ := range l_props {
				k_prop := float64(k_props[key])
				l_prop := float64(l_props[key])

				numerator += (k_prop - k_props_avg) * (l_prop - l_props_avg)

				sum_k += math.Pow((k_prop - k_props_avg), 2)
				sum_l += math.Pow((l_prop - l_props_avg), 2)
			}
			denominator := math.Sqrt(sum_k * sum_l)

			result[k_name][l_name] = numerator / denominator
		}

	}
	return result
}

func calcPropsAvg(props map[string]float64) float64 {
	var sum float64 = 0
	for _, value := range props {
		sum += value
	}
	return sum / float64(len(props))
}
