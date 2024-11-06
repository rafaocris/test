package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/qdrant/go-client/qdrant"
)

func main() {
	host := flag.String("host", "localhost", "Qdrant server host")
	port := flag.Int("port", 6334, "Qdrant server port")
	// Optional
	apiKey := flag.String("apikey", "", "API key for Qdrant (optional)")
	useTLS := flag.Bool("tls", false, "Basic TLS option")
	flag.Parse()

	// Create new client
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   *host, // Can be omitted, default is "localhost"
		Port:   *port, // Can be omitted, default is 6334
		APIKey: *apiKey,
		UseTLS: *useTLS,
		// APIKey: "<API_KEY>",
		// TLSConfig: &tls.Config{},
		// GrpcOptions: []grpc.DialOption{},
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()
	// Get a context for a minute
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	// Execute health check
	healthCheckResult, err := client.HealthCheck(ctx)
	if err != nil {
		log.Fatalf("Could not get health: %v", err)
	}
	log.Printf("Qdrant version: %s", healthCheckResult.GetVersion())

	collectionList, err := client.ListCollections(ctx)
	if err != nil {
		log.Printf("collections could not be retrieved, %v", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "Collection\tShardID\tpoints_count\t")
	fmt.Fprintln(w, "-----------------------------------")

	for _, c := range collectionList {
		clusterRequest := qdrant.CollectionClusterInfoRequest{
			CollectionName: c,
		}
		clusterInfo, err := client.GetCollectionsClient().CollectionClusterInfo(ctx, &clusterRequest)
		if err != nil {
			log.Fatal("could not get collection cluster info ", err)
		}
		for _, shard := range clusterInfo.GetLocalShards() {
			fmt.Fprintf(w, "%s\t%d\t\t\t\t\t\t\t%d\t\n", c, shard.ShardId, shard.PointsCount)
		}
		// shardCount := clusterInfo.LocalShards
	}
	w.Flush()
}
