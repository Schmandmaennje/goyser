package goyser

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/schmandmaennje/goyser/pb"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestMain(m *testing.M) {
	_, filename, _, _ := runtime.Caller(0)
	godotenv.Load(filepath.Join(filepath.Dir(filename), "..", "..", "..", "goyser", ".env"))
	os.Exit(m.Run())
}

func Test_GeyserClient(t *testing.T) {
	ctx := context.Background()

	rpcAddr, ok := os.LookupEnv("GEYSER_RPC")
	if !assert.True(t, ok, "getting GEYSER_RPC from .env") {
		t.FailNow()
	}

	if !assert.NotEqualf(t, "", rpcAddr, "GEYSER_RPC shouldn't be equal to [%s]", rpcAddr) {
		t.FailNow()
	}

	client, err := New(
		ctx,
		rpcAddr,
		nil,
	)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	defer client.Close()

	streamName := "main"
	if err = client.AddStreamClient(ctx, streamName); err != nil {
		t.Fatal(err)
	}

	stream := client.GetStreamClient(streamName)

	if err = stream.SubscribeAccounts("accounts", &geyser_pb.SubscribeRequestFilterAccounts{
		Account: []string{"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"},
	}); err != nil {
		t.Fatal(err)
	}

	for out := range stream.Ch {
		log.Printf("%+v", out)
		return
	}
}
