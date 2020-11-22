package tests

import (
	"context"
	"github.com/maglink/products-fetcher/internal/server"
	"github.com/maglink/products-fetcher/pkg/messages"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
	"testing"
)

func TestAll(t *testing.T) {
	cfg, _ := server.ReadConfig("../configs/default.yaml")

	//Run static server with CSV
	go func() {
		fs := http.FileServer(http.Dir("./"))
		http.Handle("/", fs)

		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	//Run gRPC server
	go func() {
		srv := server.New(cfg, context.Background())
		srv.Run()
	}()

	conn, err := grpc.Dial("localhost:"+strconv.Itoa(cfg.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := messages.NewProductsFetcherClient(conn)

	r, err := client.Fetch(context.Background(), &messages.FetchRequest{Url: "http://localhost:3000/products_example.csv"})
	if err != nil {
		t.Fatalf("could not fetch: %v", err)
	}
	assert.Equal(t, messages.Status_OK, r.GetStatus())

	r2, err := client.List(context.Background(), &messages.ListRequest{
		Limit:  10,
		Offset: 10,
		Order: []*messages.ListRequest_OrderOptions{
			{Field: messages.ListRequest_OrderOptions_NAME, Direction: messages.ListRequest_OrderOptions_ASC},
			{Field: messages.ListRequest_OrderOptions_PRICE, Direction: messages.ListRequest_OrderOptions_DESC},
		},
	})
	if err != nil {
		t.Fatalf("could not get list: %v", err)
	}
	assert.Equal(t, messages.Status_OK, r2.GetStatus())
	assert.Equal(t, 10, len(r2.List))
	assert.Equal(t, true, r2.List[0].Name < r2.List[1].Name)
}
