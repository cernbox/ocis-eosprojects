package service

import (
	"context"

	"github.com/cernbox/ocis-eosprojects/pkg/proto/v0"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"go.opencensus.io/trace"
)

// NewTracing returns a service that instruments traces.
func NewTracing(next EosProjects) EosProjects {
	return tracing{
		next: next,
	}
}

type tracing struct {
	next EosProjects
}

// Greet implements the EosProjects interface.
func (t tracing) GetProjects(user *userpb.User) []*proto.Project {
	_, span := trace.StartSpan(context.Background(), "EosProjects.GetProjects")
	defer span.End()

	span.Annotate([]trace.Attribute{}, "Execute EosProjects.GetProjects handler")

	return t.next.GetProjects(user)
}
