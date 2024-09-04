package main

import (
	"fmt"
	"log/slog"

	"github.com/LLIEPJIOK/forum/internal/application/forum"
)

func main() {
	if err := forum.Start(); err != nil {
		slog.Default().Error(fmt.Sprintf("cannot start program: %s", err))
	}
}
