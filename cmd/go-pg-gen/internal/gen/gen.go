package gen

import (
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

func runE(_ *cobra.Command, _ []string) error {
	g := gen.NewGenerator(gen.Config{
		OutPath:          dst,
		ModelPkgPath:     "model",
		FieldNullable:    true,
		FieldCoverable:   false,
		FieldSignable:    true,
		FieldWithTypeTag: true,
	})

	db, _ := gorm.Open(postgres.Open(dsn))
	defer func() {
		if conn, _ := db.DB(); conn != nil {
			_ = conn.Close()
		}
	}()

	g.UseDB(db)

	// float types to `github.com/shopspring/decimal` [decimal.Decimal]
	g.WithDataTypeMap(map[string]func(gorm.ColumnType) string{
		"numeric": floatToDecimal,
		"decimal": floatToDecimal,
		"real":    floatToDecimal,
		"float4":  floatToDecimal,
		"float8":  floatToDecimal,
	})

	tableMetadata := g.GenerateModel(
		table,
		gen.FieldGORMTagReg("^.*$", func(tag field.GormTag) field.GormTag {
			tag.Remove("default")
			tag.Remove("comment")
			return tag
		}))
	if tableMetadata == nil {
		Log.Error().Msgf("table metadata is nil, cannot generate table %s", table)
		os.Exit(1)
	}

	if jsonCamelCase {
		for _, v := range tableMetadata.Fields {
			if slices.Contains([]string{"int64", "uint64"}, v.Type) &&
				strings.HasSuffix(strings.ToLower(v.ColumnName), "id") {
				v.Tag.Set("json", snakeToCamel(v.ColumnName)+",string")
				continue
			}

			v.Tag.Set("json", snakeToCamel(v.ColumnName))
		}
	}

	g.Execute()

	return nil
}

func floatToDecimal(gorm.ColumnType) string {
	return "decimal.Decimal"
}

func snakeToCamel(s string) string {
	words := strings.Split(s, "_")
	title := cases.Title(language.English)
	for i, word := range words {
		if i > 0 {
			words[i] = title.String(word)
		} else {
			words[i] = strings.ToLower(word)
		}
	}
	return strings.Join(words, "")
}
