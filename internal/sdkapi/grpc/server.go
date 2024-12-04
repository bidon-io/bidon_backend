package grpcserver

import (
	"context"
	pb "github.com/bidon-io/bidon-backend/pkg/proto/bidon/v1"
	v3 "github.com/bidon-io/bidon-backend/pkg/proto/com/iabtechlab/openrtb/v3"
)

type Server struct {
	pb.UnimplementedBiddingServiceServer
}

func (s *Server) Bid(ctx context.Context, o *v3.Openrtb) (*v3.Openrtb, error) {
	//req := o.GetRequest()
	// parse context
	return &v3.Openrtb{}, nil
}
