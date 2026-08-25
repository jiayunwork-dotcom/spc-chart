package rule

import "fmt"

func CheckRule1(cz ChartZones) []Violation {
	var violations []Violation
	for i, z := range cz.Zones {
		if z.ZScore > 3 || z.ZScore < -3 {
			violations = append(violations, Violation{
				Rule:     Rule1,
				Index:    i,
				StartIdx: i,
				EndIdx:   i,
				Message:  fmt.Sprintf("point %d at z=%.2f exceeds 3-sigma", i, z.ZScore),
			})
		}
	}
	return violations
}

func CheckRule2(cz ChartZones) []Violation {
	var violations []Violation
	n := len(cz.Zones)
	if n < 9 {
		return nil
	}

	for i := 8; i < n; i++ {
		allAbove := true
		allBelow := true
		for j := i - 8; j <= i; j++ {
			if cz.Zones[j].ZScore <= 0 {
				allAbove = false
			}
			if cz.Zones[j].ZScore >= 0 {
				allBelow = false
			}
		}
		if allAbove || allBelow {
			side := "above"
			if allBelow {
				side = "below"
			}
			violations = append(violations, Violation{
				Rule:     Rule2,
				Index:    i,
				StartIdx: i - 8,
				EndIdx:   i,
				Message:  fmt.Sprintf("points %d-%d: 9 consecutive %s center line", i-8, i, side),
			})
		}
	}
	return violations
}

func CheckRule3(cz ChartZones) []Violation {
	var violations []Violation
	n := len(cz.Zones)
	if n < 6 {
		return nil
	}

	for i := 5; i < n; i++ {
		allIncreasing := true
		allDecreasing := true
		for j := i - 4; j <= i; j++ {
			if cz.Zones[j].Value <= cz.Zones[j-1].Value {
				allIncreasing = false
			}
			if cz.Zones[j].Value >= cz.Zones[j-1].Value {
				allDecreasing = false
			}
		}
		if allIncreasing || allDecreasing {
			dir := "increasing"
			if allDecreasing {
				dir = "decreasing"
			}
			violations = append(violations, Violation{
				Rule:     Rule3,
				Index:    i,
				StartIdx: i - 5,
				EndIdx:   i,
				Message:  fmt.Sprintf("points %d-%d: 6 consecutive %s", i-5, i, dir),
			})
		}
	}
	return violations
}

func CheckRule4(cz ChartZones) []Violation {
	var violations []Violation
	n := len(cz.Zones)
	if n < 14 {
		return nil
	}

	for i := 13; i < n; i++ {
		alternating := true
		for j := i - 11; j <= i; j++ {
			diff := cz.Zones[j].Value - cz.Zones[j-1].Value
			prevDiff := cz.Zones[j-1].Value - cz.Zones[j-2].Value
			if diff == 0 || prevDiff == 0 || (diff > 0 && prevDiff > 0) || (diff < 0 && prevDiff < 0) {
				alternating = false
				break
			}
		}
		if alternating {
			violations = append(violations, Violation{
				Rule:     Rule4,
				Index:    i,
				StartIdx: i - 13,
				EndIdx:   i,
				Message:  fmt.Sprintf("points %d-%d: 14 consecutive alternating", i-13, i),
			})
		}
	}
	return violations
}

func CheckRule5(cz ChartZones) []Violation {
	var violations []Violation
	n := len(cz.Zones)
	if n < 3 {
		return nil
	}

	for i := 2; i < n; i++ {
		countAbove2 := 0
		for j := i - 2; j <= i; j++ {
			if cz.Zones[j].ZScore > 2 {
				countAbove2++
			}
		}
		if countAbove2 >= 2 {
			violations = append(violations, Violation{
				Rule:     Rule5,
				Index:    i,
				StartIdx: i - 2,
				EndIdx:   i,
				Message:  fmt.Sprintf("points %d-%d: %d/3 beyond +2-sigma", i-2, i, countAbove2),
			})
			continue
		}

		countBelow2 := 0
		for j := i - 2; j <= i; j++ {
			if cz.Zones[j].ZScore < -2 {
				countBelow2++
			}
		}
		if countBelow2 >= 2 {
			violations = append(violations, Violation{
				Rule:     Rule5,
				Index:    i,
				StartIdx: i - 2,
				EndIdx:   i,
				Message:  fmt.Sprintf("points %d-%d: %d/3 beyond -2-sigma", i-2, i, countBelow2),
			})
		}
	}
	return violations
}

func CheckRule6(cz ChartZones) []Violation {
	var violations []Violation
	n := len(cz.Zones)
	if n < 5 {
		return nil
	}

	for i := 4; i < n; i++ {
		countAbove1 := 0
		for j := i - 4; j <= i; j++ {
			if cz.Zones[j].ZScore > 1 {
				countAbove1++
			}
		}
		if countAbove1 >= 4 {
			violations = append(violations, Violation{
				Rule:     Rule6,
				Index:    i,
				StartIdx: i - 4,
				EndIdx:   i,
				Message:  fmt.Sprintf("points %d-%d: %d/5 beyond +1-sigma", i-4, i, countAbove1),
			})
			continue
		}

		countBelow1 := 0
		for j := i - 4; j <= i; j++ {
			if cz.Zones[j].ZScore < -1 {
				countBelow1++
			}
		}
		if countBelow1 >= 4 {
			violations = append(violations, Violation{
				Rule:     Rule6,
				Index:    i,
				StartIdx: i - 4,
				EndIdx:   i,
				Message:  fmt.Sprintf("points %d-%d: %d/5 beyond -1-sigma", i-4, i, countBelow1),
			})
		}
	}
	return violations
}

func CheckRule7(cz ChartZones) []Violation {
	var violations []Violation
	n := len(cz.Zones)
	if n < 15 {
		return nil
	}

	for i := 14; i < n; i++ {
		allWithin := true
		for j := i - 14; j <= i; j++ {
			if cz.Zones[j].ZScore > 1 || cz.Zones[j].ZScore < -1 {
				allWithin = false
				break
			}
		}
		if allWithin {
			violations = append(violations, Violation{
				Rule:     Rule7,
				Index:    i,
				StartIdx: i - 14,
				EndIdx:   i,
				Message:  fmt.Sprintf("points %d-%d: 15 consecutive within 1-sigma (stratification)", i-14, i),
			})
		}
	}
	return violations
}

func CheckRule8(cz ChartZones) []Violation {
	var violations []Violation
	n := len(cz.Zones)
	if n < 8 {
		return nil
	}

	for i := 7; i < n; i++ {
		allBeyond1 := true
		hasAbove := false
		hasBelow := false
		for j := i - 7; j <= i; j++ {
			z := cz.Zones[j].ZScore
			if z >= -1 && z <= 1 {
				allBeyond1 = false
				break
			}
			if z > 1 {
				hasAbove = true
			}
			if z < -1 {
				hasBelow = true
			}
		}
		if allBeyond1 && hasAbove && hasBelow {
			violations = append(violations, Violation{
				Rule:     Rule8,
				Index:    i,
				StartIdx: i - 7,
				EndIdx:   i,
				Message:  fmt.Sprintf("points %d-%d: 8 beyond 1-sigma both sides (mixture)", i-7, i),
			})
		}
	}
	return violations
}
