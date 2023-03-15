package host

import (
	"github.com/gotemplates/core/runtime"
	runtime2 "github.com/gotemplates/host/runtime"
	"github.com/gotemplates/timeseries-ingress/pkg/resource"
)

const (
	databaseUrlKey       = "database-url"
	pingPathKey          = "ping-path"
	postgresUrnKey       = "postgres-urn"
	postgresPgxsqlUriKey = "postgres-pgxsql-uri"

	hostConfigFmt = "fs/host/config_{env}.txt"
)

type configMap struct {
	err error
	m   map[string]string
}

var config configMap

func init() {
	config.m, config.err = resource.ReadMap(runtime2.EnvExpansion(hostConfigFmt))
}

func (c configMap) Validate() []error {
	return runtime.ValidateMap(c.m, c.err, databaseUrlKey, pingPathKey, postgresUrnKey, postgresPgxsqlUriKey)
}

func (c configMap) DatabaseUrl() string {
	return c.m[databaseUrlKey]
}

func (c configMap) PingPath() string {
	return c.m[pingPathKey]
}

func (c configMap) PostgresUrn() string {
	return c.m[postgresUrnKey]
}

func (c configMap) PostgresPgxsqlUri() string {
	return c.m[postgresPgxsqlUriKey]
}
