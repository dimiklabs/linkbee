package main

import (
	"context"

	"github.com/shafikshaon/shortlink/cmd"
)

func main() {
	ctx := context.Background()
	cmd.Execute(ctx)
}
