module github.com/bidon-io/bidon-backend

go 1.23

require (
	github.com/Masterminds/semver/v3 v3.2.1
	github.com/alexedwards/scs/gormstore v0.0.0-20240316134038-7e11d57e8885
	github.com/alexedwards/scs/v2 v2.8.0
	github.com/bool64/cache v0.4.6
	github.com/bwmarrin/snowflake v0.3.0
	github.com/cenkalti/backoff/v4 v4.2.1
	github.com/getkin/kin-openapi v0.125.0
	github.com/getsentry/sentry-go v0.25.0
	github.com/getsentry/sentry-go/otel v0.25.0
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/go-playground/validator/v10 v10.16.0
	github.com/go-redis/cache/v9 v9.0.0
	github.com/go-redis/redismock/v9 v9.2.0
	github.com/gofrs/uuid/v5 v5.0.0
	github.com/golang-jwt/jwt/v5 v5.1.0
	github.com/google/go-cmp v0.6.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus v1.0.1
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.2.0
	github.com/jackc/pgx/v5 v5.7.1
	github.com/joho/godotenv v1.5.1
	github.com/jszwec/csvutil v1.8.0
	github.com/labstack/echo-contrib v0.15.0
	github.com/labstack/echo-jwt/v4 v4.2.0
	github.com/labstack/echo/v4 v4.11.4
	github.com/lib/pq v1.10.9
	github.com/oapi-codegen/oapi-codegen/v2 v2.3.0
	github.com/oapi-codegen/runtime v1.1.1
	github.com/oschwald/maxminddb-golang v1.12.0
	github.com/prebid/openrtb/v19 v19.0.0
	github.com/pressly/goose/v3 v3.16.0
	github.com/prometheus/client_golang v1.20.5
	github.com/redis/go-redis/v9 v9.3.0
	github.com/shopspring/decimal v1.3.1
	github.com/spazzymoto/echo-scs-session v1.0.0
	github.com/twmb/franz-go v1.18.0
	github.com/uptrace/opentelemetry-go-extra/otelgorm v0.2.3
	go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho v0.46.1
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.58.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.46.1
	go.opentelemetry.io/otel v1.33.0
	go.opentelemetry.io/otel/exporters/prometheus v0.54.0
	go.opentelemetry.io/otel/metric v1.33.0
	go.opentelemetry.io/otel/sdk v1.32.0
	go.opentelemetry.io/otel/sdk/metric v1.32.0
	go.opentelemetry.io/otel/trace v1.33.0
	go.uber.org/goleak v1.3.0
	go.uber.org/zap v1.26.0
	golang.org/x/crypto v0.31.0
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa
	google.golang.org/grpc v1.68.1
	google.golang.org/protobuf v1.35.2
	gorm.io/datatypes v1.2.4
	gorm.io/driver/postgres v1.5.10
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
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/invopop/yaml v0.3.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.60.1 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/sethvargo/go-retry v0.2.4 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.9.0 // indirect
	github.com/uptrace/opentelemetry-go-extra/otelsql v0.2.3 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/vmihailenco/go-tinylfu v0.2.2 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/mod v0.22.0 // indirect
	golang.org/x/net v0.32.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/time v0.8.0 // indirect
	golang.org/x/tools v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241209162323-e6fa225c2576 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/mysql v1.5.7 // indirect
	gorm.io/hints v1.1.2 // indirect
	gorm.io/plugin/dbresolver v1.5.3 // indirect
)
