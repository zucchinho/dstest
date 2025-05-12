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

	"cloud.google.com/go/datastore/apiv1/datastorepb"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

type mockServer struct {
	datastorepb.DatastoreServer

	Addr     string
	reqItems []reqItem
	resps    []any
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
	datastorepb.RegisterDatastoreServer(srv.Gsrv, mock)
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

		if !proto.Equal(gotReq, ri.wantReq) {
			diff := cmp.Diff(gotReq, ri.wantReq, protocmp.Transform())
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

func (s *mockServer) Lookup(ctx context.Context, in *datastorepb.LookupRequest) (*datastorepb.LookupResponse, error) {
	res, err := s.popRPC(in)
	if err != nil {
		return nil, err
	}
	return res.(*datastorepb.LookupResponse), nil
}

func (s *mockServer) BeginTransaction(_ context.Context, in *datastorepb.BeginTransactionRequest) (*datastorepb.BeginTransactionResponse, error) {
	res, err := s.popRPC(in)
	if err != nil {
		return nil, err
	}
	return res.(*datastorepb.BeginTransactionResponse), nil
}

func (s *mockServer) Commit(_ context.Context, in *datastorepb.CommitRequest) (*datastorepb.CommitResponse, error) {
	res, err := s.popRPC(in)
	if err != nil {
		return nil, err
	}
	return res.(*datastorepb.CommitResponse), nil
}

func (s *mockServer) Rollback(_ context.Context, in *datastorepb.RollbackRequest) (*datastorepb.RollbackResponse, error) {
	res, err := s.popRPC(in)
	if err != nil {
		return nil, err
	}
	return res.(*datastorepb.RollbackResponse), nil
}

func (s *mockServer) RunQuery(_ context.Context, in *datastorepb.RunQueryRequest) (*datastorepb.RunQueryResponse, error) {
	res, err := s.popRPC(in)
	if err != nil {
		return nil, err
	}
	return res.(*datastorepb.RunQueryResponse), nil
}

func (s *mockServer) RunAggregationQuery(_ context.Context, in *datastorepb.RunAggregationQueryRequest) (*datastorepb.RunAggregationQueryResponse, error) {
	res, err := s.popRPC(in)
	if err != nil {
		return nil, err
	}
	return res.(*datastorepb.RunAggregationQueryResponse), nil
}

func (s *mockServer) AllocateIds(_ context.Context, in *datastorepb.AllocateIdsRequest) (*datastorepb.AllocateIdsResponse, error) {
	res, err := s.popRPC(in)
	if err != nil {
		return nil, err
	}
	return res.(*datastorepb.AllocateIdsResponse), nil
}

func (s *mockServer) ReserveIds(_ context.Context, in *datastorepb.ReserveIdsRequest) (*datastorepb.ReserveIdsResponse, error) {
	res, err := s.popRPC(in)
	if err != nil {
		return nil, err
	}
	return res.(*datastorepb.ReserveIdsResponse), nil
}
