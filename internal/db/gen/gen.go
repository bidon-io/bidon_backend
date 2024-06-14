package gen

import (
	"database/sql"
	"log"
	"path/filepath"
	"runtime"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

func GenerateModels(db *sql.DB) {
	g := newGenerator(db, gen.Config{
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	user := g.GenerateModel(
		"users",
		gen.FieldGORMTag("is_admin", func(tag field.GormTag) field.GormTag {
			return tag.Set("default", "false")
		}),
	)

	app := g.GenerateModel(
		"apps",
		gen.FieldRelate(field.BelongsTo, "User", user, &field.RelateConfig{}),
		gen.FieldType("platform_id", "PlatformID"),
		gen.FieldType("settings", "map[string]any"),
		gen.FieldGORMTag("settings", func(tag field.GormTag) field.GormTag {
			return tag.Set("serializer", "json")
		}),
	)

	demandSource := g.GenerateModel("demand_sources")

	demandSourceAccount := g.GenerateModel(
		"demand_source_accounts",
		gen.FieldRelate(field.BelongsTo, "DemandSource", demandSource, &field.RelateConfig{}),
		gen.FieldRelate(field.BelongsTo, "User", user, &field.RelateConfig{}),
		gen.FieldType("user_id", "int64"),
		gen.FieldRename("bidding", "IsBidding"),
		gen.FieldGORMTag("bidding", func(tag field.GormTag) field.GormTag {
			return tag.Set("default", "false")
		}),
	)

	g.GenerateModel(
		"app_demand_profiles",
		gen.FieldRelate(field.BelongsTo, "App", app, &field.RelateConfig{}),
		gen.FieldRelate(field.BelongsTo, "Account", demandSourceAccount, &field.RelateConfig{}),
		gen.FieldRelate(field.BelongsTo, "DemandSource", demandSource, &field.RelateConfig{}),
	)

	segment := g.GenerateModel(
		"segments",
		gen.FieldRelate(field.BelongsTo, "App", app, &field.RelateConfig{}),
		gen.FieldType("filters", "[]segment.Filter"),
		gen.FieldGORMTag("filters", func(tag field.GormTag) field.GormTag {
			return tag.Set("serializer", "json")
		}),
	)

	g.GenerateModel(
		"auction_configurations",
		gen.FieldRelate(field.BelongsTo, "App", app, &field.RelateConfig{}),
		gen.FieldRelate(field.BelongsTo, "Segment", segment, &field.RelateConfig{
			RelatePointer: true,
		}),
		gen.FieldType("segment_id", "*sql.NullInt64"),
		gen.FieldType("ad_type", "AdType"),
		gen.FieldType("rounds", "[]auction.RoundConfig"),
		gen.FieldGORMTag("rounds", func(tag field.GormTag) field.GormTag {
			return tag.Set("serializer", "json")
		}),
		gen.FieldGORMTag("external_win_notifications", func(tag field.GormTag) field.GormTag {
			return tag.Set("default", "false")
		}),
		gen.FieldType("demands", "pq.StringArray"),
		gen.FieldType("bidding", "pq.StringArray"),
		gen.FieldType("ad_unit_ids", "pq.Int64Array"),
		gen.FieldType("settings", "map[string]any"),
		gen.FieldGORMTag("settings", func(tag field.GormTag) field.GormTag {
			return tag.Set("serializer", "json")
		}),
	)

	g.GenerateModel("countries")

	g.GenerateModel(
		"line_items",
		gen.FieldRelate(field.BelongsTo, "App", app, &field.RelateConfig{}),
		gen.FieldRelate(field.BelongsTo, "Account", demandSourceAccount, &field.RelateConfig{}),
		gen.FieldType("ad_type", "AdType"),
		gen.FieldType("extra", "map[string]any"),
		gen.FieldGORMTag("extra", func(tag field.GormTag) field.GormTag {
			return tag.Set("serializer", "json")
		}),
		gen.FieldRename("bidding", "IsBidding"),
	)

	g.Execute()
}

var dataTypeMap = map[string]func(columnType gorm.ColumnType) (dataType string){
	"bool": func(columnType gorm.ColumnType) (dataType string) {
		if n, ok := columnType.Nullable(); ok && n {
			return "sql.NullBool"
		}
		// Even if it isn't nullable we still want pointer to be able to update with `false` in almost every case
		return "*bool"
	},
	"int4": func(columnType gorm.ColumnType) (dataType string) {
		if n, ok := columnType.Nullable(); ok && n {
			return "sql.NullInt32"
		}
		return "int32"
	},
	"int8": func(columnType gorm.ColumnType) (dataType string) {
		if n, ok := columnType.Nullable(); ok && n {
			return "sql.NullInt64"
		}
		return "int64"
	},
	"numeric": func(columnType gorm.ColumnType) (dataType string) {
		if n, ok := columnType.Nullable(); ok && n {
			return "decimal.NullDecimal"
		}
		return "decimal.Decimal"
	},
	"varchar": func(columnType gorm.ColumnType) (dataType string) {
		if n, ok := columnType.Nullable(); ok && n {
			return "sql.NullString"
		}
		return "string"
	},
	"jsonb": func(columnType gorm.ColumnType) (dataType string) {
		return "datatypes.JSON"
	},
}

func newGenerator(sqlDB *sql.DB, config gen.Config) *gen.Generator {
	g := gen.NewGenerator(config)

	g.WithDataTypeMap(dataTypeMap)

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}))
	if err != nil {
		log.Fatal("open GORM db: ", err)
	}
	g.UseDB(db)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("failed to get current file path")
	}
	g.ModelPkgPath = filepath.Join(filename, "../..")

	return g
}
