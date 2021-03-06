package internal

import (
	"net/url"
	"time"

	"github.com/pinpt/agent/v4/sdk"
	"github.com/pinpt/gitlab/internal/api"
)

func (ge *GitlabExport) exportEpics(namespace *api.Namespace, repos []*api.GitlabProjectInternal, projectUsers map[string]api.UsernameMap) (rerr error) {
	if namespace.Kind == "user" {
		return
	}

	allUsers := make(api.UsernameMap)
	for _, repo := range repos {
		users := projectUsers[repo.RefID]
		for key, value := range users {
			allUsers[key] = value
		}
	}

	return api.Paginate(ge.logger, "", ge.lastExportDate, func(log sdk.Logger, params url.Values, _ time.Time) (api.NextPage, error) {
		if ge.lastExportDateGitlabFormat != "" {
			params.Set("updated_after", ge.lastExportDateGitlabFormat)
		}
		pi, epics, err := api.EpicsPage(ge.qc, namespace, params, repos)
		if err != nil {
			return pi, err
		}
		for _, epic := range epics {

			changelogs, err := ge.fetchEpicIssueDiscussions(namespace, repos, epic, allUsers)
			if err != nil {
				return pi, err
			}

			epic.ChangeLog = changelogs

			if err := ge.pipe.Write(epic); err != nil {
				return pi, err
			}
		}
		return pi, nil
	})
}
