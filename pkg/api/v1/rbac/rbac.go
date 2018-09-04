package rbac

// This package implements a simple rule based RBAC mechanism based on resource path and action.
//
// Rules are parsed as go templates against the claims for the token that was passed in with the
// request, then matched as regex against the request path and method. If both match, the request
// is granted, if either fails, it is denied.
//
// Only 2 levels of access are supported. Admins are allowed access to everything (i.e. rules are
// not evaluated at all), users are only allowed to paths they are explicitly granted access to

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
