module github.com/leodido/maintainers-generator

go 1.14

replace (
	github.com/dgrijalva/jwt-go/v4 => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.1
	k8s.io/api => k8s.io/api v0.19.3
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.3
	k8s.io/apiserver => k8s.io/apiserver v0.19.3
	k8s.io/client-go => k8s.io/client-go v0.19.3
	k8s.io/code-generator => k8s.io/code-generator v0.19.3
)

require (
	github.com/sirupsen/logrus v1.7.0
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools v2.2.0+incompatible
	k8s.io/test-infra v0.0.0-20201215091550-bbb9e1eb2f91
)
