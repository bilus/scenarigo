package main

import (
	"github.com/bilus/scenarigo/context"
	"github.com/bilus/scenarigo/testdata/gen/pb/test"
	"google.golang.org/grpc"
)

const Protocol = "grpc"

func CreateClient(ctx *context.Context, target string) test.TestClient {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		ctx.Reporter().Fatalf("failed to create client: %s", err)
	}
	return test.NewTestClient(conn)
}
