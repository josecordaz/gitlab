package api

import (
	"fmt"
	"net/url"
	"time"

	"github.com/pinpt/agent.next/sdk"
	"github.com/pinpt/go-common/v10/datetime"
	pstrings "github.com/pinpt/go-common/v10/strings"
)

func PullRequestCommentsPage(
	qc QueryContext,
	repo *sdk.SourceCodeRepo,
	pr PullRequest,
	params url.Values) (pi PageInfo, res []*sdk.SourceCodePullRequestComment, err error) {

	sdk.LogDebug(qc.Logger, "pull request comments", "repo", repo.Name, "repo_ref_id", repo.RefID, "pr", pr.IID, "params", params)

	objectPath := pstrings.JoinURL("projects", url.QueryEscape(repo.RefID), "merge_requests", pr.IID, "notes")

	var rcomments []struct {
		ID     int64 `json:"id"`
		Author struct {
			ID string `json:"id"`
		} `json:"author"`
		Body      string    `json:"body"`
		UpdatedAt time.Time `json:"updated_at"`
		CreatedAt time.Time `json:"created_at"`
		System    bool      `json:"system"`
	}

	pi, err = qc.Request(objectPath, params, &rcomments)
	if err != nil {
		return
	}

	u, err := url.Parse(qc.BaseURL)
	if err != nil {
		return pi, res, err
	}

	repoID := sdk.NewSourceCodeRepoID(qc.CustomerID, repo.RefID, qc.RefType)
	pullRequestID := sdk.NewSourceCodePullRequestID(qc.CustomerID, pr.RefID, qc.RefType, repoID)

	for _, rcomment := range rcomments {
		if rcomment.System {
			continue
		}
		item := &sdk.SourceCodePullRequestComment{}
		item.CustomerID = qc.CustomerID
		item.RefType = qc.RefType
		item.RefID = fmt.Sprint(rcomment.ID)
		item.URL = pstrings.JoinURL(u.Scheme, "://", u.Hostname(), repo.Name, "merge_requests", pr.IID)
		datetime.ConvertToModel(rcomment.UpdatedAt, &item.UpdatedDate)

		item.RepoID = repoID
		item.PullRequestID = pullRequestID
		item.Body = rcomment.Body
		datetime.ConvertToModel(rcomment.CreatedAt, &item.CreatedDate)

		item.UserRefID = rcomment.Author.ID
		res = append(res, item)
	}

	return
}
