package main

import (
	"github.com/sirupsen/logrus"
	"k8s.io/test-infra/prow/config"
	"k8s.io/test-infra/prow/git"
	gitv2 "k8s.io/test-infra/prow/git/v2"
	"k8s.io/test-infra/prow/github"
	"k8s.io/test-infra/prow/repoowners"
)

func getMaintainers(ghClient github.Client, gitClient *git.Client, repos []Repository, dedupe bool) map[string][]string {
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

	maintainers := map[string][]string{}
	for _, v := range repos {
		approvers, err := getApprovers(ownersClient, v.Org, v.Repo, dedupe)
		if err != nil {
			logrus.WithField("organization", v.Org).WithField("repository", v.Repo).Error(err)
		}
		maintainers = mergeApprovers(maintainers, approvers)
	}
	return maintainers
}
