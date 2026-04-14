// Package calculator provides mathematical utilities for gym weights.
package calculator

import (
	"sort"
)

// PlateRequirement represents how many plates of a specific weight to add to ONE SIDE of the barbell.
type PlateRequirement struct {
	Weight float64 `json:"weight"`
	Count  int     `json:"count"`
}

// CalculatePlates determines the optimal combination of plates for one side of the barbell.
func CalculatePlates(targetWeight, barbellWeight float64, availablePlates []float64) []PlateRequirement {
	if targetWeight <= barbellWeight {
		return []PlateRequirement{}
	}

	weightPerSide := (targetWeight - barbellWeight) / 2.0

	sort.Slice(availablePlates, func(i, j int) bool {
		return availablePlates[i] > availablePlates[j]
	})

	var result []PlateRequirement
	remaining := weightPerSide

	for _, plate := range availablePlates {
		if plate <= 0 {
			continue
		}

		count := int(remaining / plate)
		if count > 0 {
			result = append(result, PlateRequirement{
				Weight: plate,
				Count:  count,
			})

			remaining -= float64(count) * plate

			if remaining < 0.01 {
				break
			}
		}
	}

	return result
}
