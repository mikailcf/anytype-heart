package inboxclient

import (
	"context"

	"github.com/anyproto/any-sync/app"
	"github.com/anyproto/any-sync/coordinator/coordinatorproto"
	anysyncinboxclient "github.com/anyproto/any-sync/coordinator/inboxclient"
	"github.com/anyproto/any-sync/util/crypto"
)

// StubInboxClient is a no-op implementation of the any-sync InboxClient interface
// used in local-only mode where network services are not available
type StubInboxClient struct{}

func NewStubInboxClient() anysyncinboxclient.InboxClient {
	return &StubInboxClient{}
}

func (s *StubInboxClient) Init(a *app.App) (err error) {
	return nil
}

func (s *StubInboxClient) Name() string {
	return anysyncinboxclient.CName
}

func (s *StubInboxClient) Run(ctx context.Context) error {
	return nil
}

func (s *StubInboxClient) Close(ctx context.Context) (err error) {
	return nil
}

func (s *StubInboxClient) SetMessageReceiver(receiver anysyncinboxclient.MessageReceiver) error {
	return nil
}

func (s *StubInboxClient) InboxFetch(ctx context.Context, offset string) (messages []*coordinatorproto.InboxMessage, hasMore bool, err error) {
	return nil, false, nil
}

func (s *StubInboxClient) InboxAddMessage(ctx context.Context, receiverPubKey crypto.PubKey, message *coordinatorproto.InboxMessage) (err error) {
	return nil
}
