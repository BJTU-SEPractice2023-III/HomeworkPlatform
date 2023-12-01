package bootstrap

import "flag"

var (
	Dev              bool
	Sqlite           bool
	SqliteInMemEmpty bool
	GenData          bool
	GenDataOverwrite bool
	Mysql            bool
	Test             bool
)

func InitFlag() {
	flag.BoolVar(&Dev, "dev", false, "Whether to use dev server")
	flag.BoolVar(&Sqlite, "sqlite", false, "Whether to use sqlite")
	flag.BoolVar(&SqliteInMemEmpty, "sqlite-mem-empty", false, "Whether to use sqlite in memory (empty data)")
	flag.BoolVar(&GenData, "gen", false, "Whether to generate data to database(not overwrite)")
	flag.BoolVar(&GenDataOverwrite, "gen-overwrite", false, "Whether to generate data to database(overwrite)")
	flag.BoolVar(&Mysql, "mysql", false, "Whether to use mysql")
	flag.BoolVar(&Test, "test", false, "Whether to use test mode")
	flag.Parse()
}
