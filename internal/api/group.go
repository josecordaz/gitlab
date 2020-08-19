package api

import (
	"bytes"
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	"github.com/pinpt/agent.next/sdk"
)

// Group internal group
type Group struct {
	ID                            string
	Name                          string
	FullPath                      string
	ValidTier                     bool
	MarkedToCreateProjectWebHooks bool
	Visibility                    string
	AvatarURL                     string
}

// GroupsAll all groups
func GroupsAll(qc QueryContext) (allGroups []*Group, err error) {
	err = Paginate(qc.Logger, "", time.Time{}, func(log sdk.Logger, paginationParams url.Values, t time.Time) (np NextPage, _ error) {
		paginationParams.Set("top_level_only", "true")

		pi, groups, err := groups(qc, paginationParams)
		if err != nil {
			return pi, err
		}
		allGroups = append(allGroups, groups...)
		return pi, nil
	})
	return
}

type rawGroup struct {
	ID                 int64           `json:"id"`
	Name               string          `json:"name"`
	FullPath           string          `json:"full_path"`
	MarkedForDeletring json.RawMessage `json:"marked_for_deletion"`
	Visibility         string          `json:"visibility"`
	AvatarURL          string          `json:"avatar_url"`
}

func (g *rawGroup) reset() {
	g.ID = 0
	g.Name = ""
	g.FullPath = ""
	g.MarkedForDeletring = []byte("")
	g.Visibility = ""
	g.AvatarURL = ""
}

// Groups fetch groups
func groups(qc QueryContext, params url.Values) (np NextPage, groups []*Group, err error) {

	sdk.LogDebug(qc.Logger, "groups request", "params", sdk.Stringify(params))

	objectPath := "groups"

	var rawGroups []json.RawMessage

	np, err = qc.Get(objectPath, params, &rawGroups)
	if err != nil {
		return
	}

	var group rawGroup

	for _, g := range rawGroups {
		err = json.Unmarshal(g, &group)
		if err != nil {
			return
		}
		groups = append(groups, &Group{
			ID:         strconv.FormatInt(group.ID, 10),
			Name:       group.Name,
			FullPath:   group.FullPath,
			Visibility: group.Visibility,
			ValidTier:  isValidTier(g),
			AvatarURL:  group.AvatarURL,
		})
		group.reset()
	}

	return
}

func isValidTier(raw []byte) bool {
	return bytes.Contains(raw, []byte("marked_for_deletion"))
}

func GroupUsers(qc QueryContext, group *Group, userId string) (u *GitlabUser, err error) {

	sdk.LogDebug(qc.Logger, "group user access level", "group_name", group.Name, "group_id", group.ID, "user_id", userId)

	objectPath := sdk.JoinURL("groups", group.ID, "members", userId)

	_, err = qc.Get(objectPath, nil, &u)
	if err != nil {
		return
	}

	u.StrID = strconv.FormatInt(u.ID, 10)

	return
}

func GroupProjectsCount(qc QueryContext, group *Group) (int, error) {

	sdk.LogDebug(qc.Logger, "group projects count", "group_name", group.Name, "group_id", group.ID)

	params := url.Values{}
	params.Set("with_projects", "true")

	objectPath := sdk.JoinURL("groups", group.ID)

	var rr struct {
		Projects []json.RawMessage `json:"projects"`
	}

	_, err := qc.Get(objectPath, nil, &rr)
	if err != nil {
		return 0, err
	}

	return len(rr.Projects), nil
}
