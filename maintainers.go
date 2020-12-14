package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"k8s.io/test-infra/prow/config"
	"k8s.io/test-infra/prow/git"
	gitv2 "k8s.io/test-infra/prow/git/v2"
	"k8s.io/test-infra/prow/github"
	"k8s.io/test-infra/prow/repoowners"
)

// Maintainer is the struct representing the info about a maintainer.
type Maintainer struct {
	Name     string
	GitHub   string
	Company  string
	Projects []string
}

// Maintainers represents a list of maintainers.
type Maintainers []Maintainer

// Encode outputs the receiving maintainers list in YAML format.
func (m Maintainers) Encode() (string, error) {
	var b bytes.Buffer
	enc := yaml.NewEncoder(&b)
	err := enc.Encode(m)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func getMaintainers(ghClient github.Client, gitClient *git.Client, opts *Options) (Maintainers, error) {
	// Understand on which repository(ies) to act
	repos := []Repository{}
	if opts.repo == "" {
		var err error
		repos, err = getRepositories(ghClient, opts.org)
		if err != nil {
			return nil, err
		}
	} else {
		repos = append(repos, Repository{
			Org:  opts.org,
			Repo: opts.repo,
		})
	}

	// todo > how to change these 3 for every repo?
	enableMDYAMLFunc := func(org, repo string) bool {
		return false
	}
	skipCollaboratorsFunc := func(org, repo string) bool {
		return false
	}
	ownersDirBlacklistFunc := func() config.OwnersDirBlacklist {
		return config.OwnersDirBlacklist{}
	}

	// Get OWNERS client
	ownersClient := repoowners.NewClient(gitv2.ClientFactoryFrom(gitClient), ghClient, enableMDYAMLFunc, skipCollaboratorsFunc, ownersDirBlacklistFunc)
	// Get approvers from OWNER files for every repo
	maintainers := map[string][]string{}
	for _, v := range repos {
		approvers, err := getApprovers(ownersClient, v.Org, v.Repo, opts.dedupe)
		if err != nil {
			logrus.WithField("organization", v.Org).WithField("repository", v.Repo).Error(err)
		}
		maintainers = mergeApprovers(maintainers, approvers, opts.sort)
	}

	// Read local database
	var persons map[string]Person
	if opts.personsSupport {
		file, err := os.Open(opts.personsFile)
		defer file.Close()
		if err != nil {
			return nil, err
		}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(data, &persons)
		if err != nil {
			return nil, err
		}
	}

	// Create maintainers list
	res := Maintainers{}
	for handle, projects := range maintainers {
		host := strings.TrimSuffix(fmt.Sprintf("https://%s", opts.github.Host), "/")
		// Get the person name (or fallback to the handle)
		// Get the person company (or fallback to UNKNOWN)
		name := handle
		company := "UNKNOWN"
		person, ok := persons[handle]
		if ok {
			if len(person.Name) > 0 {
				name = person.Name
			}
			if len(person.Company) > 0 {
				company = person.Company
			}
		}

		m := Maintainer{
			Name:     name,
			GitHub:   fmt.Sprintf("%s/%s", host, handle),
			Company:  company,
			Projects: []string{},
		}
		for _, p := range projects {
			m.Projects = append(m.Projects, fmt.Sprintf("%s/%s/%s", host, opts.org, p))
		}

		res = append(res, m)
	}

	if opts.sort {
		sort.Slice(res, func(i, j int) bool {
			return res[i].Name < res[j].Name
		})
	}

	return res, nil
}
