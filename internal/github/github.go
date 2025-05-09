//go:generate go run github.com/matryer/moq@latest --pkg github_test -out logger_mock_test.go . Logger
package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/shurcooL/githubv4"

	githubREST "github.com/google/go-github/v71/github"
	"golang.org/x/oauth2"
)

// Logger formater.
type Logger interface {
	Infof(format string, opts ...interface{})
}

type defaultLogger struct{}

func (l *defaultLogger) Infof(s string, args ...interface{}) {
	fmt.Printf(
		fmt.Sprintf("[INFO] %s\n", s),
		args...,
	)
}

type localClient struct {
	url        string
	ghClient   *githubv4.Client
	token      string
	logger     Logger
	restClient *githubREST.Client
}

// Option can be passed when instantiating a client.
type Option func(s *localClient) error

// NewEnterpriseClient give a github client for use with Enterprise.
func NewEnterpriseClient(url string, token string, opts ...Option) (*localClient, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	oauthClient := oauth2.NewClient(context.Background(), src)
	lc := &localClient{
		url:        url,
		ghClient:   githubv4.NewEnterpriseClient(url, oauthClient),
		token:      token,
		logger:     &defaultLogger{},
		restClient: githubREST.NewClient(oauthClient),
	}

	for _, opt := range opts {
		if err := opt(lc); err != nil {
			return nil, fmt.Errorf("error creating client: %w", err)
		}
	}

	return lc, nil
}

// NewClient will create a new Github Client with Github's URL.
func NewClient(token string, opts ...Option) (*localClient, error) {
	return NewEnterpriseClient("https://api.github.com/graphql", token, opts...)
}

// SetVerbose will log the requests that are being made.
func SetVerbose() Option {
	return func(c *localClient) error {
		src := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: c.token},
		)
		oauthClient := oauth2.NewClient(context.Background(), src)
		oauthClient.Transport = &loggingClient{
			original: oauthClient.Transport,
			logger:   c.logger,
		}

		return nil
	}
}

// SetLogger allows to create a logger for the requests.
func SetLogger(l Logger) Option {
	return func(c *localClient) error {
		c.logger = l

		return nil
	}
}

type loggingClient struct {
	original http.RoundTripper
	logger   Logger
}

func (c *loggingClient) RoundTrip(r *http.Request) (*http.Response, error) {
	var query struct {
		Query     string
		Variables map[string]interface{}
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}

	if err := json.Unmarshal(body, &query); err != nil {
		return nil, fmt.Errorf("could not unmarshal query: %w", err)
	}

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	vars, err := json.Marshal(query.Variables)
	if err != nil {
		return nil, fmt.Errorf("could not marshal query: %w", err)
	}

	c.logger.Infof("Query: %s\nVariables: %s", query.Query, vars)

	return c.original.RoundTrip(r)
}

type PageInfo struct {
	HasNextPage bool
	EndCursor   string
}

type User struct {
	Login string
	Keys  []string
}

func (l *localClient) GetContributorKeys(ctx context.Context, owner, name string) ([]User, error) {
	users := []User{}

	var query struct {
		Repository struct {
			Name             string
			DefaultBranchRef struct {
				Target struct {
					Commit struct {
						History struct {
							PageInfo PageInfo
							Nodes    []struct {
								Author struct {
									User struct {
										Login      string
										PublicKeys struct {
											Nodes []struct {
												Key string
											}
										} `graphql:"publicKeys(first: 100, after: null)"`
									}
								}
							}
						} `graphql:"history(first: 100, after: $prCursor)"`
					} `graphql:"... on Commit"`
				}
			}
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	seen := map[string]struct{}{}

	variables := map[string]interface{}{
		"prCursor": (*githubv4.String)(nil),
		"owner":    githubv4.String(owner),
		"name":     githubv4.String(name),
	}

	for {
		if err := l.ghClient.Query(ctx, &query, variables); err != nil {
			return users, err
		}

		for _, commit := range query.Repository.DefaultBranchRef.Target.Commit.History.Nodes {
			u := commit.Author.User

			_, viewed := seen[u.Login]
			if viewed {
				continue
			}

			seen[u.Login] = struct{}{}

			keys := []string{}
			for _, key := range u.PublicKeys.Nodes {
				keys = append(keys, key.Key)
			}

			if len(keys) > 0 {
				users = append(users, User{Login: u.Login, Keys: keys})
			}
		}

		if !query.Repository.DefaultBranchRef.Target.Commit.History.PageInfo.HasNextPage {
			break
		}

		variables["prCursor"] = githubv4.String(query.Repository.DefaultBranchRef.Target.Commit.History.PageInfo.EndCursor)
	}

	return users, nil
}

func (l *localClient) GetCollaboratorKeys(ctx context.Context, owner, name string) ([]User, error) {
	users := []User{}

	var query struct {
		Repository struct {
			Collaborators struct {
				PageInfo PageInfo
				Nodes    []struct {
					Login      string
					PublicKeys struct {
						Nodes []struct {
							Key string
						}
					} `graphql:"publicKeys(first: 100, after: null)"`
				}
			} `graphql:"collaborators(first: 100, after: $prCursor)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"prCursor": (*githubv4.String)(nil),
		"owner":    githubv4.String(owner),
		"name":     githubv4.String(name),
	}

	for {
		if err := l.ghClient.Query(ctx, &query, variables); err != nil {
			return users, err
		}

		for _, collaborator := range query.Repository.Collaborators.Nodes {
			keys := []string{}
			for _, key := range collaborator.PublicKeys.Nodes {
				keys = append(keys, key.Key)
			}

			if len(keys) > 0 {
				users = append(users, User{Login: collaborator.Login, Keys: keys})
			}
		}

		if !query.Repository.Collaborators.PageInfo.HasNextPage {
			break
		}

		variables["prCursor"] = githubv4.String(query.Repository.Collaborators.PageInfo.EndCursor)
	}

	return users, nil
}

func (l *localClient) getUserKeys(ctx context.Context, author string) ([]string, error) {
	keys, _, err := l.restClient.Users.ListSSHSigningKeys(ctx, author, nil)
	if err != nil {
		return []string{}, err
	}
	k := []string{}
	for _, keys := range keys {
		k = append(k, *keys.Key)
	}
	return k, nil
}
