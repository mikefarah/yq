module github.com/mikefarah/yq

go 1.12

//replace gopkg.in/mikefarah/yaml.v2 => ../yaml

require (
	github.com/inconshreveable/mousetrap v1.0.0
	github.com/pkg/errors v0.0.0-20180311214515-816c9085562c
	github.com/spf13/pflag v0.0.0-20180601132542-3ebe029320b2
	gopkg.in/imdario/mergo.v0 v0.3.5
	gopkg.in/mikefarah/yaml.v2 v2.3.0
	gopkg.in/op/go-logging.v1 v1.0.0-20160211212156-b2cb9fa56473
	gopkg.in/spf13/cobra.v0 v0.0.3
)
