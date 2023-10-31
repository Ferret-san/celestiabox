package main

import (
	"context"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"os"

	openrpc "github.com/rollkit/celestia-openrpc"
	"github.com/rollkit/celestia-openrpc/types/blob"
	"github.com/rollkit/celestia-openrpc/types/share"
)

type DAConfig struct {
	Rpc         string `koanf:"rpc"`
	NamespaceId string `koanf:"namespace-id"`
	AuthToken   string `koanf:"auth-token"`
}

type CelestiaDA struct {
	cfg       DAConfig
	client    *openrpc.Client
	namespace share.Namespace
}

func NewCelestiaDA(cfg DAConfig) (*CelestiaDA, error) {
	daClient, err := openrpc.NewClient(context.Background(), cfg.Rpc, cfg.AuthToken)
	if err != nil {
		return nil, err
	}

	if cfg.NamespaceId == "" {
		return nil, errors.New("namespace id cannot be blank")
	}
	nsBytes, err := hex.DecodeString(cfg.NamespaceId)
	if err != nil {
		return nil, err
	}

	namespace, err := share.NewBlobNamespaceV0(nsBytes)
	if err != nil {
		return nil, err
	}

	return &CelestiaDA{
		cfg:       cfg,
		client:    daClient,
		namespace: namespace,
	}, nil
}

func (c *CelestiaDA) Store(ctx context.Context, message []byte) ([]byte, uint64, error) {
	dataBlob, err := blob.NewBlobV0(c.namespace, message)
	if err != nil {
		return nil, 0, err
	}
	commitment, err := blob.CreateCommitment(dataBlob)
	if err != nil {
		return nil, 0, err
	}
	height, err := c.client.Blob.Submit(ctx, []*blob.Blob{dataBlob}, openrpc.DefaultSubmitOptions())
	if err != nil {
		return nil, 0, err
	}
	if height == 0 {
		return nil, 0, errors.New("unexpected response code")
	}

	return commitment, height, nil
}

func (c *CelestiaDA) Read(ctx context.Context, commitment string, height uint64) ([]byte, error) {
	fmt.Println("Requesting data from Celestia", "namespace", c.cfg.NamespaceId, "commitment", commitment, "height", height)

	blob, err := c.client.Blob.Get(ctx, height, c.namespace, []byte(commitment))
	if err != nil {
		return nil, err
	}

	fmt.Println("Succesfully fetched data from Celestia", "namespace", c.cfg.NamespaceId, "height", height, "commitment", commitment)

	return blob.Data, nil
}

func readFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func writeFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

func main() {
	// Define flags
	mode := flag.String("mode", "submit", "Mode of operation: read or write")
	filename := flag.String("file", "", "Path to the file")
	commitment := flag.String("commitment", "", "commitment for the blob")
	namespace := flag.String("namespace", "000008e5f679bf7116cb", "target namespace")
	auth := flag.String("auth", "", "auth token (default is $CELESTIA_NODE_AUTH_TOKEN)")
	height := flag.Uint64("height", 0, "celestia height to fetch a blob from")

	flag.Parse()

	// Check if filename is provided
	if *auth == "" {
		fmt.Println("Please supply auth token")
		return
	}

	// Check if filename is provided
	if *filename == "" {
		fmt.Println("Please provide a filename using -file=<filename>")
		return
	}
	// Start Celestia DA
	daConfig := DAConfig{
		Rpc:         "http://localhost:26658",
		NamespaceId: *namespace,
		AuthToken:   *auth,
	}

	celestiaDA, err := NewCelestiaDA(daConfig)
	if err != nil {
		fmt.Println("Error creating Celestia client")
		return
	}

	switch *mode {
	case "submit":
		data, err := readFile(*filename)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		fmt.Printf("Successfully read %d bytes from %s\n", len(data), *filename)

		commitment, height, err := celestiaDA.Store(context.Background(), data)
		if err != nil {
			fmt.Println("Error submitting blob to Celestia")
			return
		}
		fmt.Println("Succesfully submitted blob to Celestia")
		fmt.Println("Height: ", height)
		fmt.Println("Commitment string: ", hex.EncodeToString(commitment))
	case "read":
		if *commitment == "" {
			fmt.Println("Please provide commitment using -commitment=<commitment>")
			return
		}

		if *height == 0 {
			fmt.Println("Please provide height using -height=<height>")
			return
		}

		commitmentBytes, err := hex.DecodeString(*commitment)
		if err != nil {
			fmt.Println("Error decoding hex string for commitment:", err)
			return
		}
		data, err := celestiaDA.Read(context.Background(), string(commitmentBytes), *height)
		if err != nil {
			fmt.Println("Error reading from Celestia:", err)
			return

		}
		err = writeFile(*filename, data)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		fmt.Println("File written successfully!")
	default:
		fmt.Println("Invalid mode. Please specify either 'read' or 'write'.")
	}
}
