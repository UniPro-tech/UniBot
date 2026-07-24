package main

import (
	"os"
	"path/filepath"
	"runtime"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gorm.io/gen"
)

func main() {
	//dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := os.Getenv("PG_DSN")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)

	g := gen.NewGenerator(gen.Config{
		OutPath:       filepath.Join(baseDir, "..", "..", "internal", "query"),
		Mode:          gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true, // NULL許容カラムをポインタ型として生成
	})

	g.UseDB(db)

	// Gen all tables
	g.ApplyBasic(g.GenerateAllTable()...)

	g.Execute()
}
