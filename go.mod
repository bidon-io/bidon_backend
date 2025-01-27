module github.com/bidon-io/bidon-backend

go 1.23

require (
	github.com/Masterminds/semver/v3 v3.3.1
	github.com/alexedwards/scs/gormstore v0.0.0-20240316134038-7e11d57e8885
	github.com/alexedwards/scs/v2 v2.8.0
	github.com/bool64/cache v0.4.8
	github.com/bwmarrin/snowflake v0.3.0
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/getkin/kin-openapi v0.128.0
	github.com/getsentry/sentry-go v0.31.1
	github.com/getsentry/sentry-go/echo v0.31.1
	github.com/getsentry/sentry-go/otel v0.31.1
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/go-playground/validator/v10 v10.24.0
	github.com/go-redis/cache/v9 v9.0.0
	github.com/go-redis/redismock/v9 v9.2.0
	github.com/gofrs/uuid/v5 v5.3.0
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/go-cmp v0.6.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus v1.0.1
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.2.0
	github.com/jackc/pgx/v5 v5.7.2
	github.com/joho/godotenv v1.5.1
	github.com/jszwec/csvutil v1.10.0
	github.com/labstack/echo-contrib v0.17.2
	github.com/labstack/echo-jwt/v4 v4.3.0
	github.com/labstack/echo/v4 v4.13.3
	github.com/lib/pq v1.10.9
	github.com/oapi-codegen/oapi-codegen/v2 v2.4.1
	github.com/oapi-codegen/runtime v1.1.1
	github.com/oschwald/maxminddb-golang v1.13.1
	github.com/prebid/openrtb/v19 v19.0.0
	github.com/pressly/goose/v3 v3.24.1
	github.com/prometheus/client_golang v1.20.5
	github.com/redis/go-redis/v9 v9.7.0
	github.com/shopspring/decimal v1.4.0
	github.com/spazzymoto/echo-scs-session v1.0.0
	github.com/twmb/franz-go v1.18.1
	github.com/uptrace/opentelemetry-go-extra/otelgorm v0.3.2
	go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho v0.59.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.59.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.59.0
	go.opentelemetry.io/otel v1.34.0
	go.opentelemetry.io/otel/exporters/prometheus v0.56.0
	go.opentelemetry.io/otel/metric v1.34.0
	go.opentelemetry.io/otel/sdk v1.34.0
	go.opentelemetry.io/otel/sdk/metric v1.34.0
	go.opentelemetry.io/otel/trace v1.34.0
	go.uber.org/goleak v1.3.0
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.32.0
	golang.org/x/exp v0.0.0-20250106191152-7588d65b2ba8
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.4
	gorm.io/datatypes v1.2.5
	gorm.io/driver/postgres v1.5.11
	gorm.io/gen v0.3.26
	gorm.io/gorm v1.25.12
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/invopop/yaml v0.3.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	github.com/speakeasy-api/jsonpath v0.5.2 // indirect
	github.com/speakeasy-api/openapi-overlay v0.10.0 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.9.0 // indirect
	github.com/uptrace/opentelemetry-go-extra/otelsql v0.3.2 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/vmihailenco/go-tinylfu v0.2.2 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/mod v0.22.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/time v0.9.0 // indirect
	golang.org/x/tools v0.29.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250124145028-65684f501c47 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/mysql v1.5.7 // indirect
	gorm.io/hints v1.1.2 // indirect
	gorm.io/plugin/dbresolver v1.5.3 // indirect
)
