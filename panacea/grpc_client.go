package panacea

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/tx"
	oracletypes "github.com/medibloc/panacea-core/v2/x/oracle/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GRPCClient interface {
	Close() error
	BroadcastTx(txBytes []byte) (*tx.BroadcastTxResponse, error)
	GetCdc() *codec.ProtoCodec
	GetChainID() string
	GetOraclePubKey(ctx context.Context) (*btcec.PublicKey, error)
}

var _ GRPCClient = &grpcClient{}

type grpcClient struct {
	conn    *grpc.ClientConn
	cdc     *codec.ProtoCodec
	chainID string
}

func makeInterfaceRegistry() sdk.InterfaceRegistry {
	interfaceRegistry := sdk.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	oracletypes.RegisterInterfaces(interfaceRegistry)
	return interfaceRegistry
}

func NewGRPCClient(grpcAddr, chainID string) (GRPCClient, error) {
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
		conn:    conn,
		cdc:     codec.NewProtoCodec(makeInterfaceRegistry()),
		chainID: chainID,
	}, nil
}

func (c *grpcClient) Close() error {
	log.Info("closing Panacea gRPC connection")
	return c.conn.Close()
}

func (c *grpcClient) BroadcastTx(txBytes []byte) (*tx.BroadcastTxResponse, error) {
	txClient := tx.NewServiceClient(c.conn)

	return txClient.BroadcastTx(
		context.Background(),
		&tx.BroadcastTxRequest{
			Mode:    tx.BroadcastMode_BROADCAST_MODE_BLOCK,
			TxBytes: txBytes,
		},
	)
}

func (c *grpcClient) GetCdc() *codec.ProtoCodec {
	return c.cdc
}

func (c *grpcClient) GetChainID() string {
	return c.chainID
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
