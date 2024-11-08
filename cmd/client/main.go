package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	productv1 "github.com/sekthor/grpc-streaming-example/api/product/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	ctx := context.Background()
	port := os.Getenv("PORT")

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", port), opts...)
	if err != nil {
		log.Fatal("could not dial service: %v", err)
	}
	defer conn.Close()

	client := productv1.NewProductServiceClient(conn)

	// ---------------------
	// unary rpc
	// ---------------------

	request := productv1.GetProductRequest{Id: 3}
	response, err := client.GetProduct(ctx, &request)
	if err != nil {
		log.Fatalf("cloud not get product: %v", err)
	}

	log.Println(response.Product)

	// ---------------------
	// server side streaming
	// ---------------------

	stream, err := client.GetProductList(ctx, &productv1.GetProductListRequest{})

	if err != nil {
		log.Fatalf("could not stream product list from server: %v", err)
	}

	for {
		product, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to receive streamed product: %v", err)
		}
		log.Println(product)
	}

}
