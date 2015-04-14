package github

import "github.com/crosbymichael/octokat"

type GitHub struct {
	AuthToken string
	User      string
	PR        *octokat.PullRequest
	Repo      octokat.Repo
	Content   *pullRequestContent
}

func New(token, user string, prHook *octokat.PullRequestHook) (g GitHub, err error) {
	g = GitHub{
		AuthToken: token,
		User:      user,
		PR:        prHook.PullRequest,
		Repo:      getRepo(prHook.Repo),
	}

	g.Content, err = g.getPullRequestContent(g.Repo, prHook.Number)
	if err != nil {
		return g, err
	}

	return g, nil
}

func (g GitHub) Client() *octokat.Client {
	gh := octokat.NewClient()
	gh = gh.WithToken(g.AuthToken)
	return gh
}

func getRepo(repo *octokat.Repository) octokat.Repo {
	return octokat.Repo{
		Name:     repo.Name,
		UserName: repo.Owner.Login,
	}
}
