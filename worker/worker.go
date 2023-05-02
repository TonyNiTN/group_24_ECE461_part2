package worker

import (
	"encoding/json"
	"io/ioutil"

	"github.com/19chonm/461_1_23/api"
	"github.com/19chonm/461_1_23/fileio"
	"github.com/19chonm/461_1_23/logger"
	"github.com/19chonm/461_1_23/metrics"
)

type Conts struct {
	int1 int
	int2 int
	int3 int
}

func RunTask(url string) *fileio.Rating {

	// Convert url to Github URL
	github_url, err := api.GetGithubUrl(url)
	if err != nil {
		// fmt.Println("worker: ERROR Unable to get github url ", url, " Error:", err)
		logger.DebugMsg("worker: ERROR Unable to get github url ", url, " Error:", err.Error())
	}

	license_key, err := api.GetRepoLicense(github_url)
	if err != nil {
		// fmt.Println("worker: ERROR Unable to get data for ", github_url, " License Errored:", err)
		logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, " License Errored:", err.Error())
	}

	// Get Data from Github API

	avg_lifespan, err := api.GetRepoIssueAverageLifespan(github_url)
	if err != nil {
		// fmt.Println("worker: ERROR Unable to get data for ", github_url, " AvgLifespan Errored:", err)
		logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, " AvgLifespan Errored:", err.Error())
	}

	top_recent_commits, total_recent_commits, err := api.GetRepoContributors(github_url)
	if err != nil {
		// fmt.Println("worker: ERROR Unable to get data for ", github_url, " ContributorsCommits Errored:", err)
		logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, " ContributorsCommits Errored:", err.Error())
	}
	total_prs, reviewed_prs, err := api.GetReviewFactors(github_url)
	if err != nil {
		// fmt.Println("worker: ERROR Unable to get data for ", github_url, " ScanRepo Errored:", err)
		logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, "GetReviewFactors Errored:", err.Error())
	}

	readme, err := api.GetRepoReadme(github_url)
	if err != nil {
		// fmt.Println("worker: ERROR Unable to get data for ", github_url, " ScanRepo Errored:", err)
		logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, "GetRepoReadme Errored:", err.Error())
	}
	// get repository readme

	depMap, err := api.GetRepoDependency(github_url)
	if err != nil {
		logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, " Dependency Errored:", err.Error())
	}
	//Get dependency data from github

	watchers, stargazers, totalCommits, err := api.GetCorrectnessFactors(github_url)
	if err != nil {
		// fmt.Println("worker: ERROR Unable to get data for ", github_url, " GetCorrectnessFactors Errored:", err)
		logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, " GetCorrectnessFactors Errored:", err.Error())
		return nil
	}

	// Compute scores
	correctness_score := metrics.ComputeCorrectness(watchers, stargazers, totalCommits) // no data yet
	responsiveness_score := metrics.ComputeResponsiveness(avg_lifespan)
	busfactor_score := metrics.ComputeBusFactor(top_recent_commits, total_recent_commits)
	license_score := metrics.ComputeLicenseScore(license_key)
	rampup_score := metrics.ComputeRampTimeScore(readme)
	version_score := metrics.ComputeVersionScore(depMap)
	review_score := metrics.ComputeReviewScore(int(total_prs), int(reviewed_prs))

	rampup_factor := metrics.Factor{Weight: 0.15, Value: rampup_score, AllOrNothing: false}
	correctness_factor := metrics.Factor{Weight: 0.15, Value: correctness_score, AllOrNothing: false}
	responsiveness_factor := metrics.Factor{Weight: 0.4, Value: responsiveness_score, AllOrNothing: false}
	busfactor_factor := metrics.Factor{Weight: 0.3, Value: busfactor_score, AllOrNothing: false}
	license_factor := metrics.Factor{Weight: 1.0, Value: float64(license_score), AllOrNothing: true}
	version_factor := metrics.Factor{Weight: 1.0, Value: version_score, AllOrNothing: false}
	review_factor := metrics.Factor{Weight: 1.0, Value: review_score, AllOrNothing: false}

	// Produce final rating
	factors := []metrics.Factor{rampup_factor, correctness_factor, responsiveness_factor, busfactor_factor, license_factor, version_factor, review_factor}
	r := fileio.Rating{NetScore: metrics.ComputeNetScore(factors),
		Rampup:         rampup_score,
		Url:            url,
		License:        float64(license_score),
		Busfactor:      busfactor_score,
		Responsiveness: responsiveness_score,
		Correctness:    correctness_score,
		Version:        version_score,
		Review:         review_score,
	}

	return &r
}

func ReadMapFromFile(filename string) (map[string]interface{}, error) {
	// Read file contents into byte slice
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON into map[string]interface{}
	var data map[string]interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func WriteMapToFile(data map[string]interface{}, filename string) error {
	// Convert map to JSON bytes
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Write JSON bytes to file
	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
