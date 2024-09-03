module github.com/gravitational/teleport

go 1.22.0

toolchain go1.22.4

replace github.com/codahale/hdrhistogram => github.com/HdrHistogram/hdrhistogram-go v1.1.2

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.10

replace gopkg.in/ahmetb/go-linq.v3 => github.com/ahmetb/go-linq/v3 v3.2.0

replace github.com/coreos/go-oidc => github.com/gravitational/go-oidc v0.0.1

replace github.com/docker/docker => github.com/gravitational/moby v1.4.2-0.20191008111026-2adf434ca696

replace google.golang.org/genproto => google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2
	github.com/Microsoft/go-winio v0.6.2
	github.com/aws/aws-sdk-go v1.54.16
	github.com/beevik/etree v1.4.0
	github.com/boltdb/bolt v1.3.1
	github.com/coreos/go-oidc v0.0.0-00010101000000-000000000000
	github.com/coreos/go-semver v0.3.1
	github.com/cyphar/filepath-securejoin v0.2.5
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/docker v0.0.0-00010101000000-000000000000
	github.com/fatih/color v1.17.0
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v1.3.2
	github.com/gokyle/hotp v0.0.0-20160218004637-c180d57d286b
	github.com/golang/protobuf v1.5.4
	github.com/google/gops v0.3.28
	github.com/gravitational/configure v0.0.0-20221215172404-91e9092a0318
	github.com/gravitational/form v0.0.0-20221215172421-ca521a6428ea
	github.com/gravitational/kingpin v2.1.11-0.20160205192003-785686550a08+incompatible
	github.com/gravitational/oxy v0.0.0-20231219172753-f855322f2a6c
	github.com/gravitational/roundtrip v1.0.2
	github.com/gravitational/trace v1.1.17
	github.com/gravitational/ttlmap v0.0.0-20171116003245-91fd36b9004c
	github.com/jonboulle/clockwork v0.4.0
	github.com/json-iterator/go v1.1.12
	github.com/julienschmidt/httprouter v1.3.0
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/kr/pty v1.1.1
	github.com/kylelemons/godebug v1.1.0
	github.com/mailgun/lemma v0.0.0-20170619173223-4214099fb348
	github.com/mailgun/timetools v0.0.0-20170619190023-f3a7b8ffff47
	github.com/mailgun/ttlmap v0.0.0-20170619185759-c1c17f74874f
	github.com/pborman/uuid v1.2.1
	github.com/pkg/errors v0.9.1
	github.com/pquerna/otp v1.4.0
	github.com/prometheus/client_golang v1.18.0
	github.com/russellhaering/gosaml2 v0.9.1
	github.com/russellhaering/goxmldsig v1.4.0
	github.com/sirupsen/logrus v1.9.3
	github.com/tstranex/u2f v0.0.0-20160508205855-eb799ce68da4
	github.com/vulcand/predicate v1.2.0
	go.etcd.io/etcd/api/v3 v3.5.14
	go.etcd.io/etcd/client/v3 v3.5.5
	golang.org/x/crypto v0.21.0
	golang.org/x/net v0.23.0
	golang.org/x/text v0.14.0
	google.golang.org/grpc v1.59.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.30.2
	k8s.io/apimachinery v0.30.2
	k8s.io/client-go v0.30.2
	k8s.io/klog/v2 v2.120.1
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/alecthomas/assert v1.0.0 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/coreos/pkg v0.0.0-20240122114842-bbd7aa9bf6fb // indirect
	github.com/emicklei/go-restful/v3 v3.11.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mailgun/holster v3.0.0+incompatible // indirect
	github.com/mailgun/metrics v0.0.0-20170714162148-fd99b46995bd // indirect
	github.com/mailgun/minheap v0.0.0-20170619185613-3dbe6c6bf55f // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattermost/xml-roundtrip-validator v0.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mdp/rsc v0.0.0-20160131164516-90f07065088d // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.14 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.17.0 // indirect
	golang.org/x/oauth2 v0.16.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/term v0.18.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/ahmetb/go-linq.v3 v3.0.0-00010101000000-000000000000 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools v2.2.0+incompatible // indirect
	k8s.io/kube-openapi v0.0.0-20240228011516-70dd3763d340 // indirect
	k8s.io/utils v0.0.0-20230726121419-3b25d923346b // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace google.golang.org/genproto/googleapis/rpc/status => google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1

replace google.golang.org/genproto/googleapis/rpc/code => google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1

replace google.golang.org/genproto/googleapis/rpc/errdetails => google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1

replace google.golang.org/genproto/googleapis/api/annotations => google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1

replace google.golang.org/genproto/googleapis/api/httpbody => google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/coreos/etcd/api/v3 => go.etcd.io/etcd/api/v3 v3.5.5
