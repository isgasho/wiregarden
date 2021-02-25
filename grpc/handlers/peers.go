package handlers

import (
	"context"
	"log"

	"github.com/moznion/wiregarden/grpc/messages"
	"github.com/moznion/wiregarden/internal/service"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Peers struct {
	peerService *service.Peer
	messages.UnimplementedPeersServer
}

func NewPeers(peerService *service.Peer) *Peers {
	return &Peers{
		peerService: peerService,
	}
}

func (h *Peers) GetPeers(ctx context.Context, req *messages.GetPeersRequest) (*messages.GetPeersResponse, error) {
	deviceName := req.DeviceName
	if deviceName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "device_name is a mandatory parameter, but missing")
	}

	gotPeers, err := h.peerService.GetPeers(ctx, deviceName, req.FilterPublicKeys)
	if err != nil {
		log.Printf("[error] %s", err)
		return nil, status.Error(codes.Internal, "failed to collect the peers")
	}

	peers := make([]*messages.Peer, len(gotPeers))
	for i, peer := range gotPeers {
		peers[i] = messages.ConvertFromWgctrlPeer(&peer)
	}

	return &messages.GetPeersResponse{Peers: peers}, nil
}

func (h *Peers) RegisterPeers(ctx context.Context, req *messages.RegisterPeersRequest) (*messages.RegisterPeersResponse, error) {
	deviceName := req.DeviceName
	if deviceName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "device_name is a mandatory parameter, but missing")
	}

	peers := make([]wgtypes.Peer, len(req.Peers))
	for i, reqPeer := range req.Peers {
		peer, err := reqPeer.ToWgctrlPeer()
		if err != nil {
			return nil, err
		}
		peers[i] = *peer
	}

	err := h.peerService.RegisterPeers(ctx, deviceName, peers)
	if err != nil {
		log.Printf("[error] %s", err)
		return nil, status.Error(codes.Internal, "failed to register peers")
	}

	return &messages.RegisterPeersResponse{}, nil
}

func (h *Peers) DeletePeers(context.Context, *messages.DeletePeersRequest) (*messages.DeletePeersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePeers not implemented")
}
