package github

import (
	"context"
	"os"
	"testing"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

func TestGithub(t *testing.T) {
	ctx := context.Background()
	t.Run("GetCollaboratorKeys", func(t *testing.T) {
		client, stop := newTestClient(t, "testdata/collaborators")
		defer stop()
		users, err := client.GetCollaboratorKeys(ctx, "frankywahl", "allowedSignersFile")
		if err != nil {
			t.Fatal(err)
		}

		if len(users) != 1 {
			t.Fatalf("expected 1 user, got %v", len(users))
		}

		if len(users[0].Keys) != 3 {
			t.Fatal("expected 3 keys, got ", len(users[0].Keys))
		}
	})

	t.Run("GetContributorKeys", func(t *testing.T) {
		client, stop := newTestClient(t, "testdata/contributors")
		defer stop()
		users, err := client.GetContributorKeys(ctx, "frankywahl", "allowedSignersFile")
		if err != nil {
			t.Fatal(err)
		}

		if len(users) != 1 {
			t.Fatalf("expected 1 user, got %v", len(users))
		}

		if len(users[0].Keys) != 3 {
			t.Fatal("expected 3 keys, got ", len(users[0].Keys))
		}
	})
}

func newTestClient(t *testing.T, recordLocation string) (*localClient, func() error) {
	ctx := context.Background()
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_API_TOKEN")},
	)
	oauthClient := oauth2.NewClient(ctx, src)

	r, err := recorder.NewWithOptions(&recorder.Options{
		CassetteName:  recordLocation,
		Mode:          recorder.ModeRecordOnce,
		RealTransport: oauthClient.Transport,
	})
	if err != nil {
		t.Fatal(err)
	}

	if r.Mode() != recorder.ModeRecordOnce {
		t.Fatal("Recorder should be in ModeRecordOnce")
	}
	r.AddHook(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		return nil
	}, recorder.AfterCaptureHook)

	client := &localClient{

		url:      "https://api.github.com/graphql",
		ghClient: githubv4.NewEnterpriseClient("https://api.github.com/graphql", r.GetDefaultClient()),
		token:    os.Getenv("GITHUB_API_TOKEN"),
		logger:   &defaultLogger{},
	}
	return client, r.Stop

}
