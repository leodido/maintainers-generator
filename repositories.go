package main

import (
	"fmt"

	"k8s.io/test-infra/prow/github"
)

// Repository represents the full name of a repo.
type Repository struct {
	Org    string
	Repo   string
	Branch string
}

func getRepositories(ghClient github.Client, org string) ([]Repository, error) {
	repos := []Repository{}
	// Obtain all the repositories in the organization
	res, err := ghClient.GetRepos(org, false)
	if err != nil {
		return nil, fmt.Errorf("Unable to obtain GitHub repositories")
	}
	for _, v := range res {
		// Ignore archived repositories
		if v.Archived {
			continue
		}
		// Ignore private repositories
		if v.Private {
			continue
		}
		repos = append(repos, Repository{
			Org:    org,
			Repo:   v.Name,
			Branch: v.DefaultBranch,
		})
	}
	if len(repos) == 0 {
		return nil, fmt.Errorf("Empty repository list")
	}
	return repos, nil
}
