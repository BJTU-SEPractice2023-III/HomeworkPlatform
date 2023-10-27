package bootstrap

import "flag"

var Dev bool
var Sqlite bool

func init() {
    flag.BoolVar(&Dev, "dev", false, "Whether to use dev server")
    flag.BoolVar(&Sqlite, "sqlite", false, "Whether to use sqlite")
    flag.Parse()
}
