//go:build go1.17
// +build go1.17

package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	// This function call tells the Go compiler that the interface grpc.DialOption will be used.
	// It prevents the bug that fails to method call by deleting references by the linker.
	// ref. https://github.com/bilus/scenarigo/issues/136
	grpc.Dial("", grpc.WithTransportCredentials(insecure.NewCredentials())) //nolint:errcheck
}
