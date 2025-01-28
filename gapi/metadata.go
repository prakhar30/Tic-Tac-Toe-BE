package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcUserAgentHeader = "user-agent"
	grpcAuthorityHeader = ":authority"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mdtd := &Metadata{}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// log.Printf("metadata: %v", md)
		if userAgent := md.Get(grpcUserAgentHeader); len(userAgent) > 0 {
			mdtd.UserAgent = userAgent[0]
		}
		// if clientIP := md.Get(grpcAuthorityHeader); len(clientIP) > 0 {
		// 	mdtd.ClientIP = clientIP[0]
		// }
	}

	if peerInfo, ok := peer.FromContext(ctx); ok {
		mdtd.ClientIP = peerInfo.Addr.String()
	}

	return mdtd
}
