package mipscalls

import _ "embed"

//go:embed data/mips.csv
var SyscallCsv []byte

type Syscall struct {
	Id   int    `json:"id" csv:"_number"`
	Name string `json:"name" csv:"_name"`
}
