package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/leodido/maintainers-generator/pkg/version"

	"github.com/sirupsen/logrus"
	"k8s.io/test-infra/prow/config/secret"
)

func main() {
	o := NewOptions()
	// Early exit in case user wants to only know the version
	if o.version {
		fmt.Println(version.String())
		os.Exit(0)
	}
	// Validate flags
	if err := o.Validate(); err != nil {
		logrus.WithError(err).Fatal("Invalid options")
		os.Exit(1)
	}

	// Create a secrets agent
	secretsAgent := &secret.Agent{}
	if err := secretsAgent.Start([]string{o.github.TokenPath}); err != nil {
		logrus.WithError(err).Fatal("Unable to start secrets agent.")
	}

	// Create a GitHub client
	// todo > use o.github.GitHubClientWithAccessToken() ?
	ghClient, err := o.github.GitHubClient(secretsAgent, o.dryRun)
	if err != nil {
		logrus.WithError(err).Fatal("Unable to instantiate GitHub client.")
	}

	// Create a Git client
	gitClient, err := o.github.GitClient(secretsAgent, o.dryRun)
	if err != nil {
		logrus.WithError(err).Fatal("Unable to instantiate git client.")
	}

	// Obtain maintainers
	maintainers, err := getMaintainers(ghClient, gitClient, o)
	if err != nil {
		logrus.WithField("organization", o.org).Fatal(err)
	}

	// Output YAML YaY!
	out, err := maintainers.Encode()
	if err != nil {
		logrus.WithError(err).Fatal("Unable to generate YAML maintainers file.")
	}

	var output *os.File
	if o.outputFile == "stdout" || o.outputFile == "/dev/stdout" {
		output = os.Stdout
	} else {
		output, err = os.Create(o.outputFile)
		if err != nil {
			logrus.WithError(err).WithField("output", o.outputFile).Fatal("Unable to create output file.")
		}
		defer func() {
			if err := output.Close(); err != nil {
				logrus.WithError(err).WithField("output", o.outputFile).Fatal("Unable to close output file.")
			}
		}()
	}

	w := bufio.NewWriter(output)
	n, err := w.WriteString(out)
	if err != nil {
		logrus.WithError(err).WithField("output", o.outputFile).Fatal("Unable to write output file.")
	}
	w.Flush()
	if o.outputFile != "stdout" {
		logrus.WithField("no. of bytes", n).Info("Wrote output file.")
	}
}
