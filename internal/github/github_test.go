package github_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/frankywahl/allowedSignatures/internal/github"
)

func TestGithub(t *testing.T) {
	_, _ = github.NewEnterpriseClient("foo", "bar", github.SetLogger(&LoggerMock{}), github.SetVerbose())
	ctx := t.Context()

	t.Run("GetCollaboratorKeys", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expect := "Bearer anything"
			if r.Header.Get("Authorization") != expect {
				t.Fatalf("expted Authorization Header %v, got %v", expect, r.Header.Get("Authorization"))
			}

			f, err := os.Open("testdata/GetCollaboratorKeys.json")
			if err != nil {
				t.Fatal(err)
			}

			defer f.Close()

			if _, err := io.Copy(w, f); err != nil {
				t.Fatal(err)
			}
		}))

		client, err := github.NewEnterpriseClient(s.URL, "anything")
		if err != nil {
			t.Fatal(err)
		}

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
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expect := "Bearer anything"
			if r.Header.Get("Authorization") != expect {
				t.Fatalf("expted Authorization Header %v, got %v", expect, r.Header.Get("Authorization"))
			}

			f, err := os.Open("testdata/GetContributorKeys.json")
			if err != nil {
				t.Fatal(err)
			}

			defer f.Close()

			if _, err := io.Copy(w, f); err != nil {
				t.Fatal(err)
			}
		}))

		client, err := github.NewEnterpriseClient(s.URL, "anything")
		if err != nil {
			t.Fatal(err)
		}

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
