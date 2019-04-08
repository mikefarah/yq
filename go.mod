module github.com/mikefarah/yq

go 1.12

//uncomment for local testing:
//replace github.com/mikefarah/yaml => ../yaml

//uncomment for testing with github.com/udhos/yaml:
//replace github.com/mikefarah/yaml => github.com/udhos/yaml v0.0.0-20190408152634-649037afe55b

require (
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/pkg/errors v0.0.0-20180311214515-816c9085562c
	github.com/spf13/pflag v0.0.0-20180601132542-3ebe029320b2 // indirect
	gopkg.in/imdario/mergo.v0 v0.3.5
	gopkg.in/op/go-logging.v1 v1.0.0-20160211212156-b2cb9fa56473
	gopkg.in/spf13/cobra.v0 v0.0.3
)
