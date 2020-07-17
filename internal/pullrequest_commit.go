package internal

import (
	"net/url"
	"time"

	"github.com/pinpt/agent.next.gitlab/internal/api"
	"github.com/pinpt/agent.next/sdk"
)

func (ge *GitlabExport) fetchPullRequestsCommits(repo *sdk.SourceCodeRepo, pr api.PullRequest) (commits []*sdk.SourceCodePullRequestCommit, rerr error) {
	rerr = api.Paginate(ge.logger, "", time.Time{}, func(log sdk.Logger, params url.Values, t time.Time) (api.NextPage, error) {
		pi, commitsArr, err := api.PullRequestCommitsPage(ge.qc, repo, pr, params)
		if err != nil {
			return pi, err
		}

		for _, c := range commitsArr {
			commits = append(commits, c)
		}
		return pi, nil
	})

	return

}

func (ge *GitlabExport) exportPullRequestCommits(repo *sdk.SourceCodeRepo, pr api.PullRequest) (rerr error) {

	sdk.LogDebug(ge.logger, "exporting pull requests commits", "pr", pr.Identifier)

	commits, err := ge.fetchPullRequestsCommits(repo, pr)
	if err != nil {
		rerr = err
		return
	}

	setPullRequestCommits(pr.SourceCodePullRequest, commits)
	if err := ge.writePullRequestCommits(commits); err != nil {
		rerr = err
		return
	}

	return
}
