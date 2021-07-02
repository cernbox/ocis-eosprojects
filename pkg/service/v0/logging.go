package service

import (
	"fmt"
	"time"

	"github.com/cernbox/ocis-eosprojects/pkg/proto/v0"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

// NewLogging returns a service that logs messages.
func NewLogging(next EosProjects, logger log.Logger) EosProjects {
	return logging{
		next:   next,
		logger: logger,
	}
}

type logging struct {
	next   EosProjects
	logger log.Logger
}

// Greet implements the EosProjects interface.
func (l logging) GetProjects(user *userpb.User) []*proto.Project {
	start := time.Now()
	project := l.next.GetProjects(user)

	l.logger.Debug().
		Str("method", "EosProjects.GetProjects").
		Dur("duration", time.Since(start)).
		Str("username", user.Username).
		Str("groups", fmt.Sprintf("%v", user.Groups)).
		Msg("")

	return project
}
