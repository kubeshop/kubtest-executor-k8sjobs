module github.com/kubeshop/kubtest-executor-k8sjobs

go 1.16

// replace  go get github.com/kubeshop/kubtest-executor-k8sjobs/pkg/newman => ./agent

require (
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kubeshop/kubtest v0.5.14
	github.com/stretchr/testify v1.7.0
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
)
