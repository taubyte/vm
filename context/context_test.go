package context

import (
	"context"
	"testing"
	"time"

	spec "github.com/taubyte/go-specs/common"
	"gotest.tools/v3/assert"
)

var (
	projectId     = "projectId"
	applicationId = "appId"
	resourceId    = "resourceId"
	branch        = "main"
	commit        = "commit"
)

func TestContext(t *testing.T) {
	baseContext := context.Background()
	ctx, err := New(baseContext, Project(projectId), Application(applicationId), Resource(resourceId), Commit(commit))
	assert.NilError(t, err)

	assert.Equal(t, ctx.Application(), applicationId)
	assert.Equal(t, ctx.Branch(), spec.DefaultBranch)
	assert.Equal(t, ctx.Commit(), commit)
	assert.Equal(t, ctx.Project(), projectId)
	assert.Equal(t, ctx.Resource(), resourceId)

	ctx, err = New(baseContext, Branch(branch))
	assert.NilError(t, err)

	assert.Equal(t, ctx.Branch(), branch)

	err = ctx.Close()
	assert.NilError(t, err)

	select {
	case <-ctx.Context().Done():
		return
	case <-time.After(time.Second):
		t.Error("expected context cancel")
		return
	}
}
