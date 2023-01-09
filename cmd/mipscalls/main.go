package main

import (
	"errors"
	"github.com/bfu4/mipscalls"
	"github.com/bfu4/mipscalls/api"
	"github.com/gocarina/gocsv"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"strconv"
)

var ErrNotFound error = errors.New("Syscall not found")

var (
	syscalls []mipscalls.Syscall
)

func init() {
	var err error
	if err = gocsv.UnmarshalBytes(mipscalls.SyscallCsv, &syscalls); err != nil {
		panic(err)
	}
}

func GetSyscall(id int) *mipscalls.Syscall {
	for i := 0; i < len(syscalls); i++ {
		if syscalls[i].Id == id {
			return &syscalls[i]
		}
	}
	return mipscalls.SyscallEmpty
}

func GetSyscallByName(name string) *mipscalls.Syscall {
	for i := 0; i < len(syscalls); i++ {
		if syscalls[i].Name == name {
			return &syscalls[i]
		}
	}
	return mipscalls.SyscallEmpty
}

func main() {
	var err error
	if err = godotenv.Load(); err != nil {
		panic(err)
	}
	srv := api.Get()
	srv.AddRoute(api.DefineRoute("/", []api.Method{api.GET, api.POST, api.DELETE}, func(ctx *fiber.Ctx) error {
		var (
			syscall *mipscalls.Syscall
			err error
			id  int64
		)
		name := ctx.Query("name", "")
		if name != "" {
			if syscall = GetSyscallByName(name); syscall != mipscalls.SyscallEmpty {
				return ctx.JSON(syscall)
			}
		}
		number := ctx.Query("id", "-1")
		if id, err = strconv.ParseInt(number, 10, 32); err == nil {
			if id < 0 {
				return ErrNotFound
			}
			return ctx.JSON(GetSyscall(int(id)))
		}
		return ErrNotFound
	}))
	go srv.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
