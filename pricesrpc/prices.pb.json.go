// Code generated by falafel 0.9.1. DO NOT EDIT.
// source: prices.proto

// +build js

package pricesrpc

import (
	"context"

	gateway "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func RegisterPricesJSONCallbacks(registry map[string]func(ctx context.Context,
	conn *grpc.ClientConn, reqJSON string, callback func(string, error))) {

	marshaler := &gateway.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
		},
	}

	registry["pricesrpc.Prices.GetPrice"] = func(ctx context.Context,
		conn *grpc.ClientConn, reqJSON string, callback func(string, error)) {

		req := &GetPriceRequest{}
		err := marshaler.Unmarshal([]byte(reqJSON), req)
		if err != nil {
			callback("", err)
			return
		}

		client := NewPricesClient(conn)
		resp, err := client.GetPrice(ctx, req)
		if err != nil {
			callback("", err)
			return
		}

		respBytes, err := marshaler.Marshal(resp)
		if err != nil {
			callback("", err)
			return
		}
		callback(string(respBytes), nil)
	}
}