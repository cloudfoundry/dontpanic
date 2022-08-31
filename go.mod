module code.cloudfoundry.org/dontpanic

go 1.18

require (
	github.com/jessevdk/go-flags v1.5.0
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/onsi/ginkgo/v2 v2.1.6
	github.com/onsi/gomega v1.20.1
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/google/go-cmp v0.5.8 // indirect
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace golang.org/x/text => golang.org/x/text v0.3.7
