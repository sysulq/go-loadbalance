package internal

import "google.golang.org/grpc/balancer"

var (
	// EmptyDoneFunc is a empty done function
	EmptyDoneFunc = func(balancer.DoneInfo) {}
)
