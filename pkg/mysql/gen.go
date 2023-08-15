package mysql

import (
	"context"
	"gorm.io/gen"
	"strings"
)

func NewGenTool(ctx context.Context, cfg DB, tables, outPath, pkgName string) {
	var (
		db = NewMysql(ctx, cfg)
		g  = gen.NewGenerator(gen.Config{
			OutPath:           outPath,
			ModelPkgPath:      pkgName,
			FieldNullable:     true,
			FieldWithIndexTag: true,
			FieldWithTypeTag:  true,
			Mode:              1,
		})
	)
	g.WithFileNameStrategy(func(tableName string) (fileName string) {
		fileName = strings.Replace(tableName, ".gen.go", ".go", 1)
		return fileName
	})
	g.UseDB(db)
	if len(tables) != 0 {
		g.ApplyBasic(g.GenerateModel(tables))
	} else {
		g.ApplyBasic(g.GenerateAllTable())
	}

	g.Execute()
}
