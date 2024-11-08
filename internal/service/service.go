package service

import (
	"context"

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

func (p ProductService) GetProductList(*productv1.GetProductListRequest, productv1.ProductService_GetProductListServer) error {
	panic("unimplemented")
}
