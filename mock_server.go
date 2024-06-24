package dstest

// Simple mock server for validating service requests.
//
// This mockServer follows the paradigm set here:
// https://github.com/googleapis/google-cloud-go/blob/main/datastore/mock_test.go
//
// You must add new methods to this server when testing additional

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/go-cmp/cmp"
	pb "google.golang.org/genproto/googleapis/datastore/v1"
	"google.golang.org/protobuf/proto"
)

type mockServer struct {
	pb.DatastoreServer

	Addr     string
	reqItems []reqItem
	resps    []interface{}
}

type reqItem struct {
	wantReq proto.Message
	adjust  func(gotReq proto.Message)
}

func newMockServer() (*mockServer, func(), error) {
	srv, err := NewServer()
	if err != nil {
		return nil, func() {}, err
	}

	mock := &mockServer{Addr: srv.Addr}
	pb.RegisterDatastoreServer(srv.Gsrv, mock)
	srv.Start()

	return mock, func() {
		srv.Close()
		mock.reset()
	}, nil
}

// addRPC adds a (request, response) pair to the server's list of expected
// interactions. The server will compare the incoming request with wantReq
// using proto.Equal. The response can be a message or an error.
//
// For the Listen RPC, resp should be a []interface{}, where each element
// is either ListenResponse or an error.
//
// Passing nil for wantReq disables the request check.
func (s *mockServer) AddRPC(wantReq proto.Message, resp interface{}) {
	s.AddRPCAdjust(wantReq, resp, nil)
}

// addRPCAdjust is like addRPC, but accepts a function that can be used
// to tweak the requests before comparison, for example to adjust for
// randomness.
func (s *mockServer) AddRPCAdjust(wantReq proto.Message, resp interface{}, adjust func(proto.Message)) {
	s.reqItems = append(s.reqItems, reqItem{wantReq, adjust})
	s.resps = append(s.resps, resp)
}

// popRPC compares the request with the next expected (request, response) pair.
// It returns the response, or an error if the request doesn't match what
// was expected or there are no expected rpcs.
func (s *mockServer) popRPC(gotReq proto.Message) (interface{}, error) {
	if len(s.reqItems) == 0 {
		panic(fmt.Sprintf("out of RPCs, saw %v", reflect.TypeOf(gotReq)))
	}
	ri := s.reqItems[0]
	s.reqItems = s.reqItems[1:]
	if ri.wantReq != nil {
		if ri.adjust != nil {
			ri.adjust(gotReq)
		}

		gotReqString, err := proto.Marshal(gotReq)
		if err != nil {
			return nil, fmt.Errorf("mockServer: failed to marshal got request: %v", err)
		}
		wantReqString, err := proto.Marshal(ri.wantReq)
		if err != nil {
			return nil, fmt.Errorf("mockServer: failed to marshal want request: %v", err)
		}
		if !proto.Equal(gotReq, ri.wantReq) {
			diff := cmp.Diff(gotReqString, wantReqString)
			return nil, fmt.Errorf("mockServer: bad request\ngot:%T\nwant:%T\n-got\n+want:\n%s",
				gotReq,
				ri.wantReq,
				diff,
			)
		}
	}
	resp := s.resps[0]
	s.resps = s.resps[1:]
	if err, ok := resp.(error); ok {
		return nil, err
	}
	return resp, nil
}

func (s *mockServer) reset() {
	s.reqItems = nil
	s.resps = nil
}

func (s *mockServer) Lookup(ctx context.Context, in *pb.LookupRequest) (*pb.LookupResponse, error) {
	res, err := s.popRPC(in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.LookupResponse), nil
}

func (s *mockServer) BeginTransaction(_ context.Context, in *pb.BeginTransactionRequest) (*pb.BeginTransactionResponse, error) {
	res, err := s.popRPC(in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.BeginTransactionResponse), nil
}

func (s *mockServer) Commit(_ context.Context, in *pb.CommitRequest) (*pb.CommitResponse, error) {
	res, err := s.popRPC(in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.CommitResponse), nil
}

func (s *mockServer) Rollback(_ context.Context, in *pb.RollbackRequest) (*pb.RollbackResponse, error) {
	res, err := s.popRPC(in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.RollbackResponse), nil
}

func (s *mockServer) RunQuery(_ context.Context, in *pb.RunQueryRequest) (*pb.RunQueryResponse, error) {
	res, err := s.popRPC(in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.RunQueryResponse), nil
}
