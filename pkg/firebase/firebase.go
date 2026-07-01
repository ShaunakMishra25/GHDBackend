package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/messaging"
	"firebase.google.com/go/v4/storage"
	"google.golang.org/api/option"
)

type Client struct {
	Auth     *auth.Client
	Messaging *messaging.Client
	Storage  *storage.Client
	App      *firebase.App
}

func NewClient(ctx context.Context, credentialsJSON string, storageBucket string) (*Client, error) {
	opt := option.WithCredentialsJSON([]byte(credentialsJSON))

	app, err := firebase.NewApp(ctx, &firebase.Config{
		StorageBucket: storageBucket,
	}, opt)
	if err != nil {
		return nil, fmt.Errorf("init firebase app: %w", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("init firebase auth: %w", err)
	}

	messagingClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("init firebase messaging: %w", err)
	}

	storageClient, err := app.Storage(ctx)
	if err != nil {
		return nil, fmt.Errorf("init firebase storage: %w", err)
	}

	return &Client{
		Auth:      authClient,
		Messaging: messagingClient,
		Storage:   storageClient,
		App:       app,
	}, nil
}
