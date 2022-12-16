module code.cloudfoundry.org/dontpanic

go 1.18

require (
	github.com/jessevdk/go-flags v1.5.0
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/onsi/ginkgo/v2 v2.6.1
	github.com/onsi/gomega v1.24.1
	golang.org/x/sys v0.3.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	golang.org/x/net v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace golang.org/x/text => golang.org/x/text v0.3.7
