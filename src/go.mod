module github.com/dunv/webhookForwarder

go 1.13

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/dunv/connectionTools v0.0.2
	github.com/dunv/uhelpers v1.0.11
	github.com/dunv/uhttp v1.0.43
	github.com/dunv/ulog v1.0.9
	github.com/golang/protobuf v1.3.2
	github.com/google/uuid v1.1.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pelletier/go-toml v1.6.0 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.4.0
	golang.org/x/sys v0.0.0-20200124204421-9fbb57f87de9 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/grpc v1.26.0
	gopkg.in/ini.v1 v1.51.1 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

// replace github.com/dunv/connectionTools => ../../connectionTools
