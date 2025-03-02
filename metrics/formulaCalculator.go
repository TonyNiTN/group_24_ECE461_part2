package metrics

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
)

type Factor struct {
	Weight       float64
	Value        float64
	AllOrNothing bool
}

func ComputeNetScore(fs []Factor) float64 {
	var sum float64

	for _, f := range fs {
		// if f.AllOrNothing {
		// 	if f.Value == 0 {
		// 		return 0
		// 	} else {
		// 		continue
		// 	}
		// }

		sum += f.Value * f.Weight
	}

	if sum > 1 {
		sum = 1
	} else if sum < 0 {
		sum = 0
	}

	return sum
}

func ComputeVersionScore(pckJson string) float64 {
	var pckJsonMap map[string]interface{}
	json.Unmarshal([]byte(pckJson), &pckJsonMap)
	dependencies, ok := pckJsonMap["dependencies"]
	var dependenciesMap map[string]interface{}
	if ok {
		dependenciesMap = dependencies.(map[string]interface{})
	}
	pinnedDep := 0
	var depScore float64
	re, _ := regexp.Compile(`(\d+)\.(\d+)\.(\d+)`)

	if dependenciesMap != nil {
		for _, version := range dependenciesMap {
			v := fmt.Sprintf("%v", version)
			res := re.MatchString(v)
			if res {
				pinnedDep++
			}
		}
		depScore = float64(pinnedDep) / float64(len(dependenciesMap))
	} else {
		depScore = 0.0
	}

	return depScore
}

func ComputeReviewScore(all_prs int, reviewed_prs int) float64 {
	if all_prs != 0 && reviewed_prs != 0 {
		return float64(reviewed_prs) / float64(all_prs)
	}
	return 0.0
}

func ComputeRampTimeScore(readme string) float64 {
	// Compute Ramp-up time based on number of phrases found in README
	rampUpTime := 0.0
	res, _ := regexp.MatchString(`(?i)docs\b`, readme)
	if res {
		rampUpTime = rampUpTime + 0.25
	}

	res, _ = regexp.MatchString(`(?i)quick start\b`, readme)
	if res {
		rampUpTime = rampUpTime + 0.25
	}

	res, _ = regexp.MatchString(`(?i)installation\b`, readme)
	if res {
		rampUpTime = rampUpTime + 0.25
	}

	res, _ = regexp.MatchString(`(?i)example\b`, readme)
	if res {
		rampUpTime = rampUpTime + 0.25
	}

	return rampUpTime
}

func ComputeCorrectness(watchers int, stargazers int, commits int) float64 {
	// Correctness is determined by a sum three factors: watchers, stargazers, and commits
	// Each of these are calculated using a an exponential decay function to ensure that
	// the domain is from -inf to inf, while the range is still between 0 and the weight (0.117, 0.550, or 0.333)
	// As a rough benchmark, I determined the quantity of each metric to reach a certain output value.

	// Example:
	// For watchers, the weight is 0.117. To get 80% of that weight, we need the repository to have 2000 watchers
	// The result is cs = 0.117 * 0.8 = 0.0936

	var cs, vs, ms float64
	cs = 0.117 * (1 - math.Exp(-0.001*float64(watchers)))     // 2k watchers for 80%
	vs = 0.550 * (1 - math.Exp(-0.00002*float64(stargazers))) // 100k stargazers for 86%
	ms = 0.333 * (1 - math.Exp(-0.0005*float64(commits)))     // 6000 commits for 90% of this

	return cs + vs + ms
}

func ComputeResponsiveness(days float64) float64 {
	// Compute the responsiveness score based on average
	// number of days to fix bug issues
	if days < 0 {
		return 0
	}

	return math.Exp(-0.05 * float64(days))
}

func ComputeBusFactor(top int, total int) float64 {
	// Compute the Bus factor by measuring the percentage of commits
	// in the past year committed by the top three performers
	if total <= 0 {
		return 0
	}

	return 1 - (float64(top) / float64(total))
}

func ComputeLicenseScore(license string) int {
	valid_licenses := []string{
		"agpl-3.0",     // GNU Affero General Public License v3.0
		"apache-2.0",   // Apache License 2.0
		"bsd-2-clause", // FreeBSD (https://en.wikipedia.org/wiki/BSD_licenses)
		"bsd-3-clause", // Modified BSD License
		"bsl-1.0",      // Boost Software License 1.0
		"gpl-2.0",      // GNU General Public License v2.0
		"gpl-3.0",      // GNU General Public License v3.0
		"lgpl-2.1",     // GNU Lesser General Public License v2.1
		"mit",          // MIT License
		"mpl-2.0",      // Mozilla Public License 2.0
	}

	for _, l := range valid_licenses {
		if license == l {
			return 1
		}
	}

	return 0
}
