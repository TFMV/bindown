module github.com/willabides/bindown/v3

go 1.14

require (
	code.cloudfoundry.org/bytefmt v0.0.0-20200131002437-cf55d5288a48 // indirect
	github.com/alecthomas/kong v0.2.3
	github.com/andybalholm/brotli v1.0.0 // indirect
	github.com/frankban/quicktest v1.4.2 // indirect
	github.com/killa-beez/gopkgs/sets/builtins v0.0.0-20191206232703-3018f97f77a9
	github.com/mholt/archiver/v3 v3.3.0
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/pierrec/cmdflag v0.0.2 // indirect
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/schollz/progressbar v1.0.0 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/udhos/equalfile v0.3.0
	gopkg.in/yaml.v2 v2.2.4
)

replace github.com/alecthomas/kong => github.com/willabides/kong v0.2.3-0.20200313223825-65cdca836316
