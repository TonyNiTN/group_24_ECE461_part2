package worker

import (
	"fmt"

	"github.com/19chonm/461_1_23/api"
	"github.com/19chonm/461_1_23/fileio"
	"github.com/19chonm/461_1_23/logger"
	"github.com/19chonm/461_1_23/metrics"
)

func runTask(url string, woutputch chan<- fileio.WorkerOutput) {
	// fmt.Println("My job is", url)
	logger.InfoMsg("My job is", url)

	// Convert url to Github URL
	github_url, err := api.GetGithubUrl(url)
	if err != nil {
		// fmt.Println("worker: ERROR Unable to get github url ", url, " Error:", err)
		logger.DebugMsg("worker: ERROR Unable to get github url ", url, " Error:", err.Error())
		woutputch <- fileio.WorkerOutput{WorkerErr: fmt.Errorf("worker: ERROR Unable to get github url %s  Error: %s", url, err.Error())}
		return
	}

	//Get dependency data from github
	depMap, err := api.GetRepoDependency(github_url)
	if err != nil {
		logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, " Dependency Errored:", err.Error())
		woutputch <- fileio.WorkerOutput{WorkerErr: fmt.Errorf("worker: ERROR Unable to get github url %s  Dependency Errored: %s", url, err.Error())}
		return
	}

	// get repository readme
	readme, err := api.GetRepoReadme(github_url)
	if err != nil {
		// fmt.Println("worker: ERROR Unable to get data for ", github_url, " ScanRepo Errored:", err)
		logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, " ScanRepo Errored:", err.Error())
		woutputch <- fileio.WorkerOutput{WorkerErr: fmt.Errorf("worker: ERROR Unable to get data for %s  ScanRepo Errored: %s", url, err.Error())}
		return
	}

	// Get Data from Github API
	license_key, err := api.GetRepoLicense(github_url)
	if err != nil {
		// fmt.Println("worker: ERROR Unable to get data for ", github_url, " License Errored:", err)
		logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, " License Errored:", err.Error())
		woutputch <- fileio.WorkerOutput{WorkerErr: fmt.Errorf("worker: ERROR Unable to get github url %s  License Errored: %s", url, err.Error())}
		return
	}

	avg_lifespan, err := api.GetRepoIssueAverageLifespan(github_url)
	if err != nil {
		// fmt.Println("worker: ERROR Unable to get data for ", github_url, " AvgLifespan Errored:", err)
		logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, " AvgLifespan Errored:", err.Error())
		woutputch <- fileio.WorkerOutput{WorkerErr: fmt.Errorf("worker: ERROR Unable to get data for %s  AvgLifespan Errored: %s", url, err.Error())}
		return
	}

	top_recent_commits, total_recent_commits, err := api.GetRepoContributors(github_url)
	if err != nil {
		// fmt.Println("worker: ERROR Unable to get data for ", github_url, " ContributorsCommits Errored:", err)
		logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, " ContributorsCommits Errored:", err.Error())
		woutputch <- fileio.WorkerOutput{WorkerErr: fmt.Errorf("worker: ERROR Unable to get data for %s  ContributorsCommits Errored: %s", url, err.Error())}
		return
	}

	watchers, stargazers, totalCommits, err := api.GetCorrectnessFactors(github_url)
	if err != nil {
		// fmt.Println("worker: ERROR Unable to get data for ", github_url, " GetCorrectnessFactors Errored:", err)
		logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, " GetCorrectnessFactors Errored:", err.Error())
		woutputch <- fileio.WorkerOutput{WorkerErr: fmt.Errorf("worker: ERROR Unable to get data for %s  GetCorrectnessFactors Errored: %s", url, err.Error())}
		return
	}

	// Compute scores
	correctness_score := metrics.ComputeCorrectness(watchers, stargazers, totalCommits) // no data yet
	responsiveness_score := metrics.ComputeResponsiveness(avg_lifespan)
	busfactor_score := metrics.ComputeBusFactor(top_recent_commits, total_recent_commits)
	license_score := metrics.ComputeLicenseScore(license_key)
	rampup_score := metrics.ComputeRampTime(readme)
	version_score := metrics.ComputeVersion(depMap)

	rampup_factor := metrics.Factor{Weight: 0.15, Value: rampup_score, AllOrNothing: false}
	correctness_factor := metrics.Factor{Weight: 0.15, Value: correctness_score, AllOrNothing: false}
	responsiveness_factor := metrics.Factor{Weight: 0.4, Value: responsiveness_score, AllOrNothing: false}
	busfactor_factor := metrics.Factor{Weight: 0.3, Value: busfactor_score, AllOrNothing: false}
	license_factor := metrics.Factor{Weight: 1.0, Value: float64(license_score), AllOrNothing: true}
	version_factor := metrics.Factor{Weight: 1.0, Value: version_score, AllOrNothing: false}

	// Produce final rating
	factors := []metrics.Factor{rampup_factor, correctness_factor, responsiveness_factor, busfactor_factor, license_factor, version_factor}
	r := fileio.Rating{NetScore: metrics.ComputeNetScore(factors),
		Rampup:         rampup_score,
		Url:            url,
		License:        float64(license_score),
		Busfactor:      busfactor_score,
		Responsiveness: responsiveness_score,
		Correctness:    correctness_score,
		Version:        version_score,
	}
	woutputch <- fileio.WorkerOutput{WorkerRating: r, WorkerErr: nil} // Send rating to rating channel to be sorted
}
