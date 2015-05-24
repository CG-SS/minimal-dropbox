package rest

import (
	"github.com/stretchr/testify/assert"
	"net/http"

	"context"
	"testing"
	"time"
)

func TestNopServerCloses(t *testing.T) {
	server := newNopServer()
	go server.Start()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := server.Stop(ctx)
	assert.NoError(t, err)

	select {
	case serverErr := <-server.ErrChan():
		assert.Equal(t, serverErr, http.ErrServerClosed)
	case <-time.After(1 * time.Second):
		t.FailNow()
	}
}
