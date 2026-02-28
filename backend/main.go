package main

import (
	"context"

	"github.com/shafikshaon/linkbee/cmd"
)

func main() {
	ctx := context.Background()
	cmd.Execute(ctx)
}
