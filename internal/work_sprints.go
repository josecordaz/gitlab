package internal

import (
	"net/url"
	"time"

	"github.com/pinpt/agent.next.gitlab/internal/api"
	"github.com/pinpt/agent.next/sdk"
)

func (ge *GitlabExport) exportRepoSprints(project *sdk.SourceCodeRepo) error {
	return api.Paginate(ge.logger, "", time.Time{}, func(log sdk.Logger, params url.Values, t time.Time) (pi api.NextPage, rerr error) {
		pi, sprints, err := api.SprintsPage(ge.qc, project, params)
		if err != nil {
			return pi, err
		}
		for _, s := range sprints {
			s.IntegrationInstanceID = ge.integrationInstanceID
			if err := ge.pipe.Write(s); err != nil {
				return
			}
		}
		return
	})
}
