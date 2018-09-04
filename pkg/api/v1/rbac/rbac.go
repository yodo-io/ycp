/*
Package rbac implements simple rule engine for RBAC based on a request path and action.

A set of rules is matched against the current request matching request path and method against
allowed values defined by the rule. If any rule matches, the request is granted, if no rule matches,
the request is denied.

A rule consists of a request path and an action. The path is parsed as a go template against
the claim found in the JWT token that was transmitted with the request. If no claim was found, access
is denied.

The parsed path template & action in the rule are matched as regular expressions against the request
URLs path component and HTTP method respectively.

	// JWT contained this claim
	const claim := auth.Claim {
		UserID: 1
	}

	// Path and action are matched as regex
	// Additionally, path is parsed as template against claim
	rule := &rule {
		path: "/vd+/users/{{.UserID}}"
		action: "GET|PATCH"
	}

	// The following requests would match:
	// GET /v1/users/1
	// PATCH /v1/users/1
	// PATCH /v2/users/1

	// These requests would not match:
	// GET /v1/users/2
	// DELETE /v1/users/1
	// POST /v1/users


Currently only two role-based access levels are supported. Users with RoleAdmin are always allowed access
(i.e. no rules are evaluated), users with RoleUser need to have at least one matching rule.

As of now, rules are harcoded in this package - however, the implementation would allow to easily serialize
them in any structured document format (JSON, YAML) and keep them in a database or config file.
*/
package rbac

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"text/template"

	"github.com/yodo-io/ycp/pkg/api/v1/auth"
	"github.com/yodo-io/ycp/pkg/model"

	"github.com/gin-gonic/gin"
	"github.com/yodo-io/ycp/pkg/api"
)

type rule struct {
	Path   string
	Action string
}

// We should store these rules in a config file or database, for simplicity they are hardcoded here
var allow = []rule{
	{Path: `/v\d+/catalog`, Action: `GET`},
	{Path: `/v\d+/resources/{{.UserID}}`, Action: `.*`},
	{Path: `/v\d+/quotas/{{.UserID}}`, Action: `GET`},
	{Path: `/v\d+/users/{{.UserID}}`, Action: `.*`},
}

// Middleware creates new RBAC middleware
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// No claims, no do
		o, hasClaims := c.Get("claims")
		if !hasClaims {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrStr("Unauthorized"))
		}
		cl, ok := o.(auth.Claims)
		if !ok {
			api.Fatal(c, errors.New("Invalid claim type"))
		}

		// Admin can do anything
		if cl.Role == model.RoleAdmin {
			return
		}

		// First rule that matches wins, if none matches, deny
		// Evaluate path as go template to allow for "user can access their own stuff" type rules
		// Match both as simple regexes
		for _, r := range allow {
			pm, err := pathMatcher(&cl, r)
			if err != nil {
				api.Fatal(c, err)
			}
			ac, err := actionMatcher(r)
			if err != nil {
				api.Fatal(c, err)
			}
			if pm.Match([]byte(c.Request.URL.Path)) && ac.Match([]byte(c.Request.Method)) {
				return // pass
			}
		}

		// No matching rules, no do
		c.AbortWithStatusJSON(http.StatusForbidden, api.ErrStr("Forbidden"))
	}
}

// Render path as template from rule, then compile into a regex
func pathMatcher(cl *auth.Claims, r rule) (*regexp.Regexp, error) {
	b := &strings.Builder{}
	tpl, err := template.New("").Parse(r.Path)
	if err != nil {
		return nil, err
	}
	if err := tpl.Execute(b, cl); err != nil {
		return nil, err
	}
	re, err := regexp.Compile(b.String())
	if err != nil {
		return nil, err
	}
	return re, nil
}

// Actions are simply converted into regexes, no template parsing
func actionMatcher(r rule) (*regexp.Regexp, error) {
	re, err := regexp.Compile(r.Action)
	if err != nil {
		return nil, err
	}
	return re, nil
}
