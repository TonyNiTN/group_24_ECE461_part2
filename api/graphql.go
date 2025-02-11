/*
Gets factors that are used to build the rating for correctness of modules via
GitHub's GraphQL API. getCorrectnessFactors function returns factors in the
order of watchers, stargazers, totalCommits. A data type in the same form as
the query structure is required to convert string to json. From json, the data
is returned.

--Use of GitHub token needs to be changed
*/

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CorrectnessFactors struct {
	Data struct {
		Repository struct {
			StargazerCount int64
			Watchers       struct {
				TotalCount int64
			}
			DefaultBranchRef struct {
				Target struct {
					History struct {
						TotalCount int64
					}
				}
			}
		}
	}
}

type TotalPRs struct {
	Data struct {
		Repository struct {
			Item1 struct {
				TotalCount int64 `json:"totalCount"`
			} `json:"item1"`
			Item2 struct {
				TotalCount int64 `json:"totalCount"`
			} `json:"item2"`
		} `json:"repository"`
	}
}

func buildTotalPRsQuery(ownerName string, repoName string, states string) (query map[string]string) {
	var totalPRsQuery = map[string]string{
		"query": `
		{
			repository(owner:` + `"` + ownerName + `", name:` + `"` + repoName + `") { 
				item1: pullRequests(first: 0, states: [OPEN, CLOSED, MERGED]) {
					totalCount
				}
				item2: pullRequests(first: 0, states: [MERGED]) {
					totalCount
				}
			}
		}`,
	}

	return totalPRsQuery
}

func buildCorrectnessQuery(ownerName string, repoName string) (query map[string]string) {
	var correctnessQuery = map[string]string{
		"query": `
		{
			repository(owner:` + `"` + ownerName + `", name:` + `"` + repoName + `") { 
				stargazerCount
				watchers {
					totalCount
				}
				defaultBranchRef {
					target {
						... on Commit {
							history(first:0) {
								totalCount
							}
						}
					}
				}
			}
		}`,
	}

	return correctnessQuery
}

func GetCorrectnessFactors(url string) (watchers int, stargazers int, totalCommits int, err error) {
	ownerName, repoName, token, err := ValidateInput(url)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("GetCorrectnessFactors: Error on validate input")
	}

	query := buildCorrectnessQuery(ownerName, repoName)

	jsonValue, _ := json.Marshal(query)
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(jsonValue))
	req.Header.Add("Authorization", "Bearer "+token)
	client := &http.Client{Timeout: time.Second * 10}
	res, err := client.Do(req)
	defer res.Body.Close()

	if err != nil {
		return 0, 0, 0, fmt.Errorf("The GraphQL query failed with error %s\n", err)
	}

	var factors CorrectnessFactors
	err = json.NewDecoder(res.Body).Decode(&factors)

	if err != nil {
		return 0, 0, 0, fmt.Errorf("Reading body failed with errorr %s\n", err)
	}

	watchers = int(factors.Data.Repository.StargazerCount)
	stargazers = int(factors.Data.Repository.Watchers.TotalCount)
	totalCommits = int(factors.Data.Repository.DefaultBranchRef.Target.History.TotalCount)

	return watchers, stargazers, totalCommits, nil
}


func GetReviewFactors(url string) (all_prs int, reviewd_prs int, err error) {
	ownerName, repoName, token, err := ValidateInput(url)
	if err != nil {
		return 0, 0, fmt.Errorf("GetReviewFactors: Error on validate input")
	}

	query := buildTotalPRsQuery(ownerName, repoName, "OPEN,CLOSED,MERGED")
	jsonValue, _ := json.Marshal(query)
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(jsonValue))
	req.Header.Add("Authorization", "Bearer "+token)
	client := &http.Client{Timeout: time.Second * 10}
	res, err := client.Do(req)
	defer res.Body.Close()

	if err != nil {
		return 0, 0, fmt.Errorf("The GraphQL query failed with error %s\n", err)
	}

	var factors TotalPRs
	err = json.NewDecoder(res.Body).Decode(&factors)

	if err != nil {
		return 0, 0, fmt.Errorf("Reading body failed with errorr %s\n", err)
	}
	all_prs = int(factors.Data.Repository.Item1.TotalCount)
	reviewd_prs = int(factors.Data.Repository.Item2.TotalCount)

	return all_prs, reviewd_prs, nil
}
