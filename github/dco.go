package github

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/crosbymichael/octokat"
)

func (g GitHub) DcoVerified(prHook *octokat.PullRequestHook) (bool, error) {
	// we only want the prs that are opened/synchronized
	if !prHook.IsOpened() && !prHook.IsSynchronize() {
		return false, nil
	}

	// we only want apply labels
	// to opened pull requests
	var labels []string

	//check if it's a proposal
	isProposal := strings.Contains(strings.ToLower(g.PR.Title), "proposal")
	switch {
	case isProposal:
		labels = []string{"status/1-needs-design-review"}
	case g.Content.IsDocsOnly():
		labels = []string{"status/3-needs-docs-review"}
	default:
		labels = []string{"status/0-needs-triage"}
	}

	// add labels if there are any
	// only add labels to new PRs not sync
	if len(labels) > 0 && prHook.IsOpened() {
		log.Debugf("Adding labels %#v to pr %d", labels, prHook.Number)

		if err := g.addLabel(g.Repo, prHook.Number, labels...); err != nil {
			return false, err
		}

		log.Infof("Added labels %#v to pr %d", labels, prHook.Number)
	}

	var verified bool

	if g.Content.CommitsSigned() {
		if err := g.toggleLabels(g.Repo, prHook.Number, "dco/no", "dco/yes"); err != nil {
			return false, err
		}

		if err := g.removeComment(g.Repo, g.PR, "sign your commits", g.Content); err != nil {
			return false, err
		}

		if err := g.successStatus(g.Repo, g.PR.Head.Sha, "docker/dco-signed", "All commits signed"); err != nil {
			return false, err
		}

		verified = true
	} else {
		if err := g.toggleLabels(g.Repo, prHook.Number, "dco/yes", "dco/no"); err != nil {
			return false, err
		}

		if err := g.addDCOUnsignedComment(g.Repo, g.PR, g.Content); err != nil {
			return false, err
		}

		if err := g.failureStatus(g.Repo, g.PR.Head.Sha, "docker/dco-signed", "Some commits without signature", "https://github.com/docker/docker/blob/master/CONTRIBUTING.md#sign-your-work"); err != nil {
			return false, err
		}
	}

	return verified, nil
}
