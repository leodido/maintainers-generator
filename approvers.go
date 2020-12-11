package main

import (
	"fmt"
	"reflect"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"k8s.io/test-infra/prow/repoowners"
)

func traverseProcessMap(v reflect.Value, current string, out *map[string][]string) {
	if v.Kind() == reflect.Map {
		logrus.WithField("function", "traverseProcessMap").Debugln("call")
		for _, key := range v.MapKeys() {
			logrus.WithField("key", key.Interface()).WithField("type", key.Kind()).Debugln("found key")
			content := v.MapIndex(key)

			switch key.Kind() {
			case reflect.String:
				switch content.Kind() {
				case reflect.Struct:
					username := key.Interface().(string)
					if _, ok := (*out)[username]; !ok {
						logrus.WithField("handle", username).WithField("project", current).Debugln("create entry")
						(*out)[username] = []string{current}
					} else {
						logrus.WithField("handle", username).WithField("project", current).Debugln("update entry")
						(*out)[username] = append((*out)[username], current)
					}
					break
				case reflect.Map:
					current := key.Interface().(string)
					logrus.WithField("directory", current).Debugln("detected folder to append to next handles")
					traverseProcessMap(content, current, out)
					break
				}
				break
			case reflect.Ptr:
				switch content.Kind() {
				case reflect.Map:
					traverseProcessMap(content, current, out)
					break
				}
				break
			}
		}
	}
}

func processApprovers(v reflect.Value, out *map[string][]string) error {
	if out == nil {
		return fmt.Errorf("missing output for processing approvers")
	}
	project := ""
	traverseProcessMap(v, project, out) // v is `map[string]map[*regexp.Regexp]sets.String`

	return nil
}

func getApprovers(ownersClient *repoowners.Client, org, repo string, dedupe bool) error {
	owners, err := ownersClient.LoadRepoOwners(org, repo, "master")
	if err != nil {
		logrus.WithError(err).WithField("organization", org).WithField("repository", repo).Fatal("Unable to fetch OWNERS.")
	}

	approvers := getUnexportedValue(reflect.ValueOf(owners).Elem().FieldByName("approvers"))

	spew.Dump(approvers)
	result := map[string][]string{}
	err = processApprovers(approvers, &result)
	if err != nil {
		return err
	}

	for k, values := range result {
		newvalues := []string{}
		for _, v := range values {
			if dedupe && v == "" {
				newvalues = []string{fmt.Sprintf("%s/%s", org, repo)}
				break
			}
			if v != "" {
				v = fmt.Sprintf("/%s", v)
			}
			newvalues = append(newvalues, fmt.Sprintf("%s/%s%s", org, repo, v))
		}
		result[k] = newvalues
	}

	spew.Dump(result)

	return nil
}
