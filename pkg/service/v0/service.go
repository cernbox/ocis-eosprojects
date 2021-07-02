package service

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/juliangruber/go-intersect"

	"github.com/cernbox/ocis-eosprojects/pkg/config"
	"github.com/cernbox/ocis-eosprojects/pkg/proto/v0"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/owncloud/ocis/ocis-pkg/log"

	//mysql package requires blank import
	_ "github.com/go-sql-driver/mysql"
)

var permissionsLevel = map[string]int{
	"admins":  1,
	"writers": 2,
	"readers": 3,
}

type EosProjects interface {
	GetProjects(user *userpb.User) (project []*proto.Project)
}

// New returns a new instance of Service
func NewEosProjects(dbCfg config.DB, opts ...Option) (EosProjects, error) {
	options := newOptions(opts...)

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Name))
	if err != nil {
		panic(err)
	}

	p := EosProjectsmpl{
		log:     options.Logger,
		db:      db,
		dbTable: dbCfg.Table,
	}

	return p, nil
}

// BasicGreeter implements the Greeter interface
type EosProjectsmpl struct {
	log     log.Logger
	db      *sql.DB
	dbTable string
}

// Greet implements the EosProjectsHandler interface.
func (p EosProjectsmpl) GetProjects(user *userpb.User) []*proto.Project {
	groups := []string{"cernbox-project-swan-admins", "cernbox-project-swan-writers", "cernbox-project-cernbox-writers", "cernbox-project-cernbox-readers"}

	p.log.Info().
		Str("username", user.Username).
		Msg("Getting projects")

	r := regexp.MustCompile(`^cernbox-project-(?P<Name>.+)-(?P<Permissions>admins|writers|readers)\z`)

	userProjects := make(map[string]string)
	var userProjectsKeys []string

	for _, group := range groups {
		match := r.FindStringSubmatch(group)
		if match != nil {
			if userProjects[match[1]] == "" {
				userProjectsKeys = append(userProjectsKeys, match[1])
			}
			userProjects[match[1]] = getHigherPermission(userProjects[match[1]], match[2])
		}
	}

	if len(userProjectsKeys) == 0 {
		// User has no projects... lets bail
		p.log.Info().
			Str("username", user.Username).
			Str("groups", fmt.Sprintf("%v", user.Groups)).
			Msg("User has no project egroup")
		return nil
	}

	var dbProjects []string
	dbProjectsPaths := make(map[string]string)
	query := fmt.Sprintf("SELECT project_name, eos_relative_path FROM %s", p.dbTable)
	results, err := p.db.Query(query)
	if err != nil {
		p.log.Error().
			Err(err).
			Msg("Failed to get projects from DB")
		return nil
	}

	for results.Next() {
		var name string
		var path string
		err = results.Scan(&name, &path)
		if err != nil {
			p.log.Error().
				Err(err).
				Msg("Failed to map DB to variables")
			return nil
		}
		dbProjects = append(dbProjects, name)
		dbProjectsPaths[name] = path
	}

	validProjects := intersect.Simple(dbProjects, userProjectsKeys)

	var projects []*proto.Project
	for _, project := range validProjects.([]interface{}) {
		name := project.(string)
		permissions := userProjects[name]
		projects = append(projects, &proto.Project{
			Name:        name,
			Path:        fmt.Sprintf("/eos/project/%s", dbProjectsPaths[name]), //Hardcoded for now...
			Permissions: permissions[:len(permissions)-1],
		})
	}

	return projects
}

func getHigherPermission(perm1, perm2 string) string {
	if perm1 == "" {
		return perm2
	}
	if permissionsLevel[perm1] < permissionsLevel[perm2] {
		return perm1
	}
	return perm2
}
