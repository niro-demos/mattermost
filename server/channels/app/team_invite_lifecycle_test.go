// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost/server/public/model"
)

func TestDeletedTeamInviteCannotBeResolvedOrJoined(t *testing.T) {
	mainHelper.Parallel(t)
	th := Setup(t).InitBasic(t)

	team := th.CreateTeam(t)
	outsider := th.CreateUser(t)

	// Control: the invitation resolves while the team is active.
	resolved, appErr := th.App.GetTeamByInviteId(team.InviteId)
	require.Nil(t, appErr)
	require.Equal(t, team.Id, resolved.Id)

	require.Nil(t, th.App.SoftDeleteTeam(team.Id))

	resolved, appErr = th.App.GetTeamByInviteId(team.InviteId)
	require.NotNil(t, appErr, "a deleted team's invitation must not resolve")
	assert.Equal(t, http.StatusNotFound, appErr.StatusCode)
	assert.Nil(t, resolved)

	joinedTeam, member, appErr := th.App.AddUserToTeamByInviteId(th.Context, team.InviteId, outsider.Id)
	require.NotNil(t, appErr, "a deleted team's invitation must not create membership")
	assert.Equal(t, http.StatusNotFound, appErr.StatusCode)
	assert.Nil(t, joinedTeam)
	assert.Nil(t, member)

	require.Nil(t, th.App.RestoreTeam(team.Id))
	_, err := th.App.Srv().Store().Team().GetMember(th.Context, team.Id, outsider.Id)
	require.Error(t, err, "restoring the team must not activate membership attempted while deleted")

	restored, appErr := th.App.GetTeamByInviteId(team.InviteId)
	require.Nil(t, appErr)
	assert.Equal(t, model.TeamOpen, restored.Type)
}
