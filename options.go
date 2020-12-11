package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"k8s.io/test-infra/pkg/flagutil"
	prowflagutil "k8s.io/test-infra/prow/flagutil"
)

// Options represents the flag for the current plugin
type Options struct {
	dryRun     bool
	github     prowflagutil.GitHubOptions
	hmacSecret string
	version    bool
	logLevel   string
	org        string
	repo       string
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

	return nil
}

// NewOptions instantiates Options from arguments
func NewOptions() *Options {
	o := Options{}
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.BoolVar(&o.dryRun, "dry-run", true, "Dry run for testing (uses API tokens but does not mutate).")
	fs.StringVar(&o.hmacSecret, "hmac", "/etc/webhook/hmac", "Path to the file containing the GitHub HMAC secret.")
	fs.BoolVar(&o.version, "version", false, "Print the version.")
	// fs.BoolVar(&o.enableMDYAML, "enable-mdyaml", false, "Whether to enable support for MD/YAML OWNERS files.")
	// fs.BoolVar(&o.skipCollabor, "skip-collaborators", false, "Whether to skip collaborators for processing the maintainers.")
	fs.StringVar(&o.logLevel, "log-level", "info", "Log level.")
	fs.StringVar(&o.org, "org", "", "The GitHub organization name.")
	fs.StringVar(&o.repo, "repo", "", "The GitHub repository name.")

	for _, group := range []flagutil.OptionGroup{&o.github} {
		group.AddFlags(fs)
	}
	fs.Parse(os.Args[1:])

	return &o
}
