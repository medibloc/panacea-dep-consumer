package panacea

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/btcsuite/btcd/btcec"
	oracletypes "github.com/medibloc/panacea-core/v2/x/oracle/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GRPCClient interface {
	GetOraclePubKey(ctx context.Context) (*btcec.PublicKey, error)
}

var _ GRPCClient = &grpcClient{}

type grpcClient struct {
	conn *grpc.ClientConn
}

func NewGRPCClient(grpcAddr string) (GRPCClient, error) {
	log.Infof("dialing to Panacea gRPC endpoint: %s", grpcAddr)

	parsedUrl, err := url.Parse(grpcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gRPC endpoint. please use absolute URL (scheme://host:port): %w", err)
	}

	var cred grpc.DialOption

	if parsedUrl.Scheme == "https" {
		cred = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))
	} else {
		cred = grpc.WithInsecure()
	}

	conn, err := grpc.Dial(parsedUrl.Host, cred)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Panacea: %w", err)
	}

	return &grpcClient{
		conn: conn,
	}, nil
}

func (c *grpcClient) Close() error {
	log.Info("closing Panacea gRPC connection")
	return c.conn.Close()
}

func (c *grpcClient) GetOraclePubKey(ctx context.Context) (*btcec.PublicKey, error) {
	client := oracletypes.NewQueryClient(c.conn)

	response, err := client.Params(ctx, &oracletypes.QueryOracleParamsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get oracle params info via grpc: %w", err)
	}

	oraclePublicKeyBytes, err := base64.StdEncoding.DecodeString(response.GetParams().OraclePublicKey)
	if err != nil {
		return nil, err
	}

	oraclePubKey, err := btcec.ParsePubKey(oraclePublicKeyBytes, btcec.S256())
	if err != nil {
		return nil, err
	}

	return oraclePubKey, nil
}
