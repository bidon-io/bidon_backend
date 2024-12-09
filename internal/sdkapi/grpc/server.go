package grpcserver

import (
	"context"
	v3 "github.com/bidon-io/bidon-backend/pkg/proto/com/iabtechlab/openrtb/v3"
	pb "github.com/bidon-io/bidon-backend/pkg/proto/org/bidon/proto/v1"
)

type Server struct {
	pb.UnimplementedBiddingServiceServer
}

func (s *Server) Bid(ctx context.Context, o *v3.Openrtb) (*v3.Openrtb, error) {
	//req := o.GetRequest()
	// parse context
	return &v3.Openrtb{}, nil
}
