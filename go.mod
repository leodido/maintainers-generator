module github.com/leodido/maintainers-generator

go 1.14

replace (
	k8s.io/api => k8s.io/api v0.19.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.0
	k8s.io/client-go => k8s.io/client-go v0.19.0
	github.com/dgrijalva/jwt-go/v4 => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
)

require (
	github.com/sirupsen/logrus v1.6.0
	gopkg.in/yaml.v2 v2.3.0
	gotest.tools v2.2.0+incompatible
	k8s.io/test-infra v0.0.0-20201210094454-9d2209384daf
)
