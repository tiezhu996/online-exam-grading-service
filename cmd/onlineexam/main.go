package main

import (
	"fmt"

	"onlineexam/internal/config"
	"onlineexam/internal/service"
	"onlineexam/internal/store"
)

func main() {
	cfg := config.Load()
	st := store.New()
	svc := service.New(st, cfg)
	_ = svc
	fmt.Println("onlineexam ready")
}
