package worker

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

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
	// fmt.Println("My job is", url)
	logger.InfoMsg("My job is", url)
	cache := make(map[string]interface{})
	if _, err := os.Stat("cache"); os.IsNotExist(err) {
		// Create cache if it does not exist
		file, err := os.Create("cache")
		if err != nil {
			logger.DebugMsg("Error creating cache")
		}
		cache["init"] = 1
		err = WriteMapToFile(cache, "cache")
		if err != nil {
			logger.DebugMsg("Error writing to cache file")
		}
		defer file.Close()
		logger.DebugMsg("Created cache")
	} else if err != nil {
		logger.DebugMsg("Error checking for the cache file.")
	} else {
		logger.DebugMsg("cache exists")
		gob.Register(Conts{})
		cache, err = ReadMapFromFile("cache")
		if err != nil {
			logger.DebugMsg("Error loading cache from file!")
		}

	}

	// Convert url to Github URL
	github_url, err := api.GetGithubUrl(url)
	if err != nil {
		// fmt.Println("worker: ERROR Unable to get github url ", url, " Error:", err)
		logger.DebugMsg("worker: ERROR Unable to get github url ", url, " Error:", err.Error())
		return nil
	}

	var license_key string
	cachedResponse, found := cache[fmt.Sprintf("%s-license", url)]
	if found {
		license_key = cachedResponse.(string)
	} else {
		license_key, err = api.GetRepoLicense(github_url)
		if err != nil {
			// fmt.Println("worker: ERROR Unable to get data for ", github_url, " License Errored:", err)
			logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, " License Errored:", err.Error())
			return nil
		}
		cache[fmt.Sprintf("%s-license", url)] = license_key

	}
	// Get Data from Github API

	var avg_lifespan float64
	cachedResponse, found = cache[fmt.Sprintf("%s-issues", url)]
	if found {
		avg_lifespan = cachedResponse.(float64)
	} else {
		avg_lifespan, err = api.GetRepoIssueAverageLifespan(github_url)
		if err != nil {
			// fmt.Println("worker: ERROR Unable to get data for ", github_url, " AvgLifespan Errored:", err)
			logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, " AvgLifespan Errored:", err.Error())
			return nil
		}
		cache[fmt.Sprintf("%s-issues", url)] = avg_lifespan

	}

	cachedResponse, found = cache[fmt.Sprintf("%s-contributors", url)]
	var top_recent_commits int
	var total_recent_commits int
	if found {
		cr := cachedResponse.(string)
		resp := strings.Fields(cr)
		top_recent_commits, _ = strconv.Atoi(resp[0])
		total_recent_commits, _ = strconv.Atoi(resp[1])
	} else {
		top_recent_commits, total_recent_commits, err = api.GetRepoContributors(github_url)
		if err != nil {
			// fmt.Println("worker: ERROR Unable to get data for ", github_url, " ContributorsCommits Errored:", err)
			logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, " ContributorsCommits Errored:", err.Error())
			return nil
		}
		cache[fmt.Sprintf("%s-contributors", url)] = strconv.Itoa(top_recent_commits) + " " + strconv.Itoa(total_recent_commits)
	}

	//Get All pull request data from github
	// total_prs, err := api.GetRepoPRs(github_url)
	// if err != nil {
	// 	// fmt.Println("worker: ERROR Unable to get data for ", github_url, " License Errored:", err)
	// 	logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, " GetRepoPRs Errored:", err.Error())
	// 	woutputch <- fileio.WorkerOutput{WorkerErr: fmt.Errorf("worker: ERROR Unable to get github url %s  Pull Requests Errored: %s", url, err.Error())}
	// 	return
	// }

	// // Get reviewed pull request data from github
	// reviewed_prs, err := api.GetReviewedPRs(github_url)
	// if err != nil {
	// 	//fmt.Println("worker: ERROR Unable to get data for ", github_url, " Reviewed Pull Requests Errored:", err)
	// 	logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, " GetReviewedPRs Errored:", err.Error())
	// 	woutputch <- fileio.WorkerOutput{WorkerErr: fmt.Errorf("worker: ERROR Unable to get github url %s  Reviewed Pull Requests Errored: %s", url, err.Error())}
	// 	return
	// }

	var total_prs int
	var reviewed_prs int
	cachedResponse, found = cache[fmt.Sprintf("%s-review", url)]
	if found {
		cr := cachedResponse.(string)
		resp := strings.Fields(cr)
		total_prs, _ = strconv.Atoi(resp[0])
		reviewed_prs, _ = strconv.Atoi(resp[1])
	} else {
		total_prs, reviewed_prs, err = api.GetReviewFactors(github_url)
		if err != nil {
			// fmt.Println("worker: ERROR Unable to get data for ", github_url, " ScanRepo Errored:", err)
			logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, "GetReviewFactors Errored:", err.Error())
			return nil
		}
		cache[fmt.Sprintf("%s-review", url)] = strconv.Itoa(int(total_prs)) + " " + strconv.Itoa(int(reviewed_prs))
	}

	var readme string
	cachedResponse, found = cache[fmt.Sprintf("%s-readme", url)]
	if found {
		readme = cachedResponse.(string)
	} else {
		readme, err = api.GetRepoReadme(github_url)
		if err != nil {
			// fmt.Println("worker: ERROR Unable to get data for ", github_url, " ScanRepo Errored:", err)
			logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, "GetRepoReadme Errored:", err.Error())
			return nil
		}
		cache[fmt.Sprintf("%s-readme", url)] = string(readme)

	}
	// get repository readme

	var depMap string
	cachedResponse, found = cache[fmt.Sprintf("%s-dependency", url)]
	if found {
		depMap = cachedResponse.(string)
	} else {
		depMap, err = api.GetRepoDependency(github_url)
		if err != nil {
			logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, " Dependency Errored:", err.Error())
			return nil
		}
		cache[fmt.Sprintf("%s-dependency", url)] = string(depMap)

	}
	//Get dependency data from github

	var watchers, stargazers, totalCommits int
	cachedResponse, found = cache[fmt.Sprintf("%s-correctness", url)]
	if found {
		cr := cachedResponse.(string)
		resp := strings.Fields(cr)
		watchers, _ = strconv.Atoi(resp[0])
		stargazers, _ = strconv.Atoi(resp[1])
		totalCommits, _ = strconv.Atoi(resp[2])
	} else {
		watchers, stargazers, totalCommits, err = api.GetCorrectnessFactors(github_url)
		if err != nil {
			// fmt.Println("worker: ERROR Unable to get data for ", github_url, " GetCorrectnessFactors Errored:", err)
			logger.DebugMsg("worker: ERROR Unable to get data for ", github_url, " GetCorrectnessFactors Errored:", err.Error())
			return nil
		}
		cache[fmt.Sprintf("%s-correctness", url)] = strconv.Itoa(int(watchers)) + " " + strconv.Itoa(int(stargazers)) + " " + strconv.Itoa(int(totalCommits))

	}

	err = WriteMapToFile(cache, "cache")
	if err != nil {
		logger.DebugMsg("worker: ERROR writing to cache file ")
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
