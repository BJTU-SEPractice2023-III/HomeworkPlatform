package bootstrap

import "flag"

var Dev bool
var Sqlite bool
var Mysql bool

func InitFlag() {
	flag.BoolVar(&Dev, "dev", false, "Whether to use dev server")
	flag.BoolVar(&Sqlite, "sqlite", false, "Whether to use sqlite")
	flag.BoolVar(&Mysql, "mysql", false, "Whether to use mysql")
	flag.Parse()
}
