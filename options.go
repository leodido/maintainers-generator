package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"k8s.io/test-infra/pkg/flagutil"
	prowflagutil "k8s.io/test-infra/prow/flagutil"
)

// Options represents the flag for the current plugin
type Options struct {
	dryRun         bool
	github         prowflagutil.GitHubOptions
	version        bool
	logLevel       string
	org            string
	repo           string
	dedupe         bool
	sort           bool
	personsFile    string
	outputFile     string
	personsSupport bool
	banner         bool
}

// Validate validates the receiving options.
func (o *Options) Validate() error {
	for _, group := range []flagutil.OptionGroup{&o.github} {
		if err := group.Validate(o.dryRun); err != nil {
			return err
		}
	}

	lvl, err := logrus.ParseLevel(o.logLevel)
	if err != nil {
		return fmt.Errorf("%s is not a valid logrus log level", o.logLevel)
	}
	logrus.SetLevel(lvl)

	if o.org == "" && o.repo == "" {
		return fmt.Errorf("specify at least a GitHub organization")
	}

	if _, err := os.Stat(o.personsFile); err == nil {
		if filepath.Ext(o.personsFile) == ".json" {
			o.personsSupport = true
		} else {
			o.personsSupport = false
			logrus.WithField("path", o.personsFile).Warn("The persons file is not JSON, disabling support")
		}
	} else if os.IsNotExist(err) {
		o.personsSupport = false
		logrus.WithField("path", o.personsFile).Warn("The persons file does not exist, disabling support")
	} else {
		// file may or may not exist
		o.personsSupport = false
		logrus.WithField("path", o.personsFile).WithError(err).Warn("The persons file may or may not exists (see error for details), disabling support")
	}

	return nil
}

// NewOptions instantiates Options from arguments
func NewOptions() *Options {
	o := Options{}
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.BoolVar(&o.dryRun, "dry-run", false, "Dry run for testing (uses API tokens but does not mutate).")
	fs.BoolVar(&o.version, "version", false, "Print the version.")
	fs.BoolVar(&o.dedupe, "dedupe", true, "Whether to dedupe or not sub-project areas for every maintainer.")
	fs.BoolVar(&o.sort, "sort", true, "Whether to sort the projects alphabetically.")
	// fs.BoolVar(&o.enableMDYAML, "enable-mdyaml", false, "Whether to enable support for MD/YAML OWNERS files.")
	// fs.BoolVar(&o.skipCollabor, "skip-collaborators", false, "Whether to skip collaborators for processing the maintainers.")
	fs.StringVar(&o.logLevel, "log-level", "info", "Log level.")
	fs.StringVar(&o.org, "org", "", "The GitHub organization name.")
	fs.StringVar(&o.repo, "repo", "", "The GitHub repository name.")
	fs.StringVar(&o.personsFile, "persons-db", "data/data.json", "The path to a JSON file containing handle => name/company mappings")
	fs.StringVar(&o.outputFile, "output", "stdout", "The path where to write the output YAML maintainers")
	fs.BoolVar(&o.banner, "banner", false, "Whether you want a header on top of the output YAML maintainers file")

	for _, group := range []flagutil.OptionGroup{&o.github} {
		group.AddFlags(fs)
	}
	fs.Parse(os.Args[1:])

	return &o
}
