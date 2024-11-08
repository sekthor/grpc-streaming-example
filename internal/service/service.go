package service

import (
	"context"
	"io"
	"time"

	productv1 "github.com/sekthor/grpc-streaming-example/api/product/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ productv1.ProductServiceServer = ProductService{}

type ProductService struct {
	productv1.UnimplementedProductServiceServer
}

func (p ProductService) GetProduct(ctx context.Context, req *productv1.GetProductRequest) (*productv1.GetProductResponse, error) {

	var response productv1.GetProductResponse

	for _, product := range products {

		if product.ID == req.Id {
			prod := productv1.Product{}
			prod.Id = product.ID
			prod.Name = product.Name

			response.Product = &prod

			return &response, nil
		}
	}

	return &response, status.Errorf(codes.NotFound, "No product with id %d", req.Id)
}

func (p ProductService) GetProductList(req *productv1.GetProductListRequest, stream productv1.ProductService_GetProductListServer) error {

	for _, product := range products {

		prod := productv1.Product{
			Id:    product.ID,
			Name:  product.Name,
			Price: product.Price,
		}

		// silly sleep to observe individual messages from the stream
		time.Sleep(time.Second)

		if err := stream.Send(&prod); err != nil {
			return err
		}
	}
	return nil
}

func (p ProductService) FillCart(stream productv1.ProductService_FillCartServer) error {

	var cart []*productv1.Product
	var sum float32 = 0

	for {
		product, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&productv1.Cart{
				Products:   cart,
				TotalPrice: sum,
			})
		}
		if err != nil {
			return err
		}

		cart = append(cart, product)
		sum += product.Price
	}
}
