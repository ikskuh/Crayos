package meta

import "flag"

var FLAG_ADDR = flag.String("addr", ":8080", "http service address")
var DEBUG_MODE = flag.Bool("debug", false, "Enables debug mode (default session + no session death)")
