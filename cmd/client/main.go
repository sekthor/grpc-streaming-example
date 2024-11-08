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

	var products []*productv1.Product
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
		products = append(products, product)
		log.Println(product)
	}

	// ---------------------
	// client side streaming
	// ---------------------

	fillCartStream, err := client.FillCart(ctx)
	if err != nil {
		log.Fatalf("could not start fill-cart stream: %v", err)
	}

	for _, product := range products {

		if err := fillCartStream.Send(product); err != nil {
			log.Fatalf("could not send product to fill-cart stream: %v", err)
		}
	}

	cart, err := fillCartStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("could not close fill-cart stream: %v", err)
	}

	log.Println(cart)
}
