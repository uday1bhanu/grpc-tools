package replay

import (
	"context"
	"fmt"
	"github.com/bradleyjkemp/grpc-tools/grpc-proxy"
	"github.com/bradleyjkemp/grpc-tools/internal"
	"github.com/bradleyjkemp/grpc-tools/internal/codec"
	"github.com/bradleyjkemp/grpc-tools/internal/marker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"time"
)

func Run(protoRoots, protoDescriptors, dumpPath, destinationOverride string, dialer grpc_proxy.ContextDialer) error {
	//	pool := internal.NewConnPool(logrus.New(), dialer)
	//
	//	dumpFile, err := os.Open(dumpPath)
	//	if err != nil {
	//		return err
	//	}
	//	var resolvers []proto_decoder.MessageResolver
	//	if protoRoots != "" {
	//		r, err := proto_decoder.NewFileResolver(strings.Split(protoRoots, ",")...)
	//		if err != nil {
	//			return err
	//		}
	//		resolvers = append(resolvers, r)
	//	}
	//	if protoDescriptors != "" {
	//		r, err := proto_decoder.NewDescriptorResolver(strings.Split(protoRoots, ",")...)
	//		if err != nil {
	//			return err
	//		}
	//		resolvers = append(resolvers, r)
	//	}
	//	encoder := proto_decoder.NewEncoder(resolvers...)
	//
	//	dumpDecoder := json.NewDecoder(dumpFile)
	//RPC:
	//	for {
	//		rpc := dump_format.RPC{}
	//		err := dumpDecoder.Decode(&rpc)
	//		if err == io.EOF {
	//			break
	//		}
	//		if err != nil {
	//			return fmt.Errorf("failed to decode dump: %s", err)
	//		}
	//
	//		conn, err := getConnection(pool, rpc.Metadata, destinationOverride)
	//		if err != nil {
	//			return fmt.Errorf("failed to connect to destination (%s): %s", destinationOverride, err)
	//		}
	//
	//		// RPC has metadata added by grpc-dump that should be removed before sending
	//		// (so that we're sending as close as possible to the original request)
	//		marker.RemoveHTTPSMarker(rpc.Metadata)
	//
	//		ctx := metadata.NewOutgoingContext(context.Background(), rpc.Metadata)
	//		streamName := rpc.StreamName()
	//		str, err := conn.NewStream(ctx, &grpc.StreamDesc{
	//			StreamName:    streamName,
	//			ServerStreams: true,
	//			ClientStreams: true,
	//		}, streamName)
	//		if err != nil {
	//			return fmt.Errorf("failed to make new stream: %v", err)
	//		}
	//
	//		fmt.Print(streamName, "...")
	//		for _, message := range rpc.Messages {
	//			msgBytes, err := encoder.Encode(streamName, message)
	//			if err != nil {
	//				return fmt.Errorf("failed to encode message: %v", err)
	//			}
	//
	//			switch message.MessageOrigin {
	//			case dump_format.ClientMessage:
	//				err := str.SendMsg(msgBytes)
	//				if err != nil {
	//					return fmt.Errorf("failed to send message: %v", err)
	//				}
	//			case dump_format.ServerMessage:
	//				var resp []byte
	//				err := str.RecvMsg(&resp)
	//				if err != nil {
	//					// TODO when do we assert on RPC errors?
	//					return fmt.Errorf("failed to recv message: %v", err)
	//				}
	//				if string(resp) != string(msgBytes) {
	//					fmt.Println("Err mismatch")
	//					continue RPC
	//				}
	//			default:
	//				return fmt.Errorf("invalid message type: %v", message.MessageOrigin)
	//			}
	//		}
	//		fmt.Println("OK")
	//	}
	return nil
}

func getConnection(pool *internal.ConnPool, md metadata.MD, destinationOverride string) (*grpc.ClientConn, error) {
	// if no destination override set then auto-detect from the metadata
	var destination = destinationOverride
	if destination == "" {
		authority := md.Get(":authority")
		if len(authority) == 0 {
			return nil, fmt.Errorf("no destination override specified and could not auto-detect from dump")
		}
		destination = authority[0]
	}

	options := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NoopCodec{})),
		grpc.WithBlock(),
	}

	if marker.IsTLSRPC(md) {
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	} else {
		options = append(options, grpc.WithInsecure())
	}

	dialCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return pool.GetClientConn(dialCtx, destination, options...)
}
