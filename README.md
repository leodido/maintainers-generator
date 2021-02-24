# maintainers-generator

> Parse OWNERS file across repositories to output YAML containing all the maintainers!

## Install

```console
go get github.com/leodido/maintainers-generator
```

## Usage

This is the CLI.

```console
Usage of ./bin/maintainers-generator:
  -banner
        Whether you want a header on top of the output YAML maintainers file
  -dedupe
        Whether to dedupe or not sub-project areas for every maintainer. (default true)
  -dry-run
        Dry run for testing (uses API tokens but does not mutate).
  -github-app-id string
        ID of the GitHub app. If set, requires --github-app-private-key-path to be set and --github-token-path to be unset.
  -github-app-private-key-path string
        Path to the private key of the github app. If set, requires --github-app-id to bet set and --github-token-path to be unset
  -github-endpoint value
        GitHub's API endpoint (may differ for enterprise). (default https://api.github.com)
  -github-graphql-endpoint string
        GitHub GraphQL API endpoint (may differ for enterprise). (default "https://api.github.com/graphql")
  -github-host string
        GitHub's default host (may differ for enterprise) (default "github.com")
  -github-token-path string
        Path to the file containing the GitHub OAuth secret.
  -log-level string
        Log level. (default "info")
  -org string
        The GitHub organization name.
  -output string
        The path where to write the output YAML maintainers (default "stdout")
  -persons-db string
        The path to a JSON file containing handle => name/company mappings (default "data/data.json")
  -repo string
        The GitHub repository name.
  -sort
        Whether to sort the projects alphabetically. (default true)
  -version
        Print the version.
```

For example, you could run:

```console
./bin/maintainers-generator --github-token-path /etc/token --org falcosecurity
```

Which will output a YAML to STDOUT like the following one:

```yaml
- name: user1
  github: https://github.com/user1
  company: UNKNOWN
  projects:
  - https://github.com/falcosecurity/falcosecurity/community
- name: user2
  github: https://github.com/user2
  company: UNKNOWN
  projects:
  - https://github.com/falcosecurity/falcosecurity/.github
  - https://github.com/falcosecurity/falcosecurity/advocacy
  - https://github.com/falcosecurity/falcosecurity/charts
  - https://github.com/falcosecurity/falcosecurity/client-go
  - https://github.com/falcosecurity/falcosecurity/client-rs
```

## TODOs

- Extensive test suite
- Complete support for the [OWNERS spec](https://github.com/kubernetes/community/blob/master/contributors/guide/owners.md)
  - Support for `no_parent_owners: true`
  - Support for OWNERS file with `filters`
  - Support for `emeritus_approvers`

---

[![Analytics](https://ga-beacon.appspot.com/UA-49657176-1/maintainers-generator?flat)](https://github.com/igrigorik/ga-beacon)
