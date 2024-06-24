package dstest

import (
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewDatastoreServer(t *testing.T) (*grpc.ClientConn, *mockServer, func()) {
	srv, cleanup, err := newMockServer()
	if err != nil {
		t.Fatal(err)
	}
	conn, err := grpc.Dial(
		srv.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())
	if err != nil {
		t.Fatal(err)
	}
	return conn, srv, func() {
		conn.Close()
		cleanup()
	}
}
