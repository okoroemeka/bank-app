package gapi

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"log"
)

type Metadata struct {
	UserAgent string
	ClientIp  string
}

const (
	xForwardedHost       = "x-forwarded-host"
	grpcGatewayUserAgent = "grpcgateway-user-agent"
	userAgent            = "user-agent"
)

func (server *Server) extractMetadata(context context.Context) *Metadata {
	mtdt := &Metadata{}
	if md, ok := metadata.FromIncomingContext(context); ok {
		log.Printf("metadata: %+v\n", md)
		if ua, ok := md[grpcGatewayUserAgent]; ok {
			mtdt.UserAgent = ua[0]
		}
		if ip, ok := md[xForwardedHost]; ok {
			mtdt.ClientIp = ip[0]
		}
		if ua, ok := md[userAgent]; ok {
			mtdt.UserAgent = ua[0]
		}
	}
	if pr, ok := peer.FromContext(context); ok {
		mtdt.ClientIp = pr.Addr.String()
	}
	return mtdt
}
