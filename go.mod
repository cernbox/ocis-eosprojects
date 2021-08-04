module github.com/cernbox/ocis-eosprojects

go 1.16

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.7.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.2
	github.com/asim/go-micro/v3 v3.5.1-0.20210217182006-0f0ace1a44a9
	github.com/cs3org/go-cs3apis v0.0.0-20210726064751-4e567a2c8df1
	github.com/cs3org/reva v1.10.1-0.20210730095301-fcb7a30a44a6
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/juliangruber/go-intersect v1.0.0
	github.com/micro/cli/v2 v2.1.2
	github.com/oklog/run v1.1.0
	github.com/openzipkin/zipkin-go v0.2.5
	github.com/owncloud/ocis v1.9.1-0.20210804055949-69dfdeac0d9c
	github.com/prometheus/client_golang v1.10.0
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/thejerf/suture/v4 v4.0.1
	go.opencensus.io v0.23.0
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
	google.golang.org/protobuf v1.27.1
)

replace (
	github.com/crewjam/saml => github.com/crewjam/saml v0.4.5
	github.com/owncloud/ocis => /root/ocis
	github.com/owncloud/ocis/ocis-pkg => github.com/owncloud/ocis v1.9.1-0.20210804055949-69dfdeac0d9c
	go.etcd.io/etcd/api/v3 => go.etcd.io/etcd/api/v3 v3.0.0-20210204162551-dae29bb719dd
	go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3 v3.0.0-20210204162551-dae29bb719dd
)
