package authz

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/go-faster/errors"
	"gopkg.in/yaml.v3"
)

type EngineConfig struct {
	ModelPath  string
	PolicyPath string
	RoutesPath string
}

type Engine struct {
	enforcer *casbin.Enforcer
	routes   []compiledRule
}

// Enforcer exposes underlying casbin enforcer for introspection/testing.
func (e *Engine) Enforcer() *casbin.Enforcer {
	return e.enforcer
}

type Permission struct {
	Domain   string
	Resource string
	Action   string
}

func (p Permission) String() string {
	return fmt.Sprintf("%s:%s:%s", p.Domain, p.Resource, p.Action)
}

func parsePermission(raw string) (Permission, error) {
	parts := strings.Split(raw, ":")
	if len(parts) != 3 {
		return Permission{}, fmt.Errorf("invalid permission %q, want domain:resource:action", raw)
	}
	return Permission{
		Domain:   strings.TrimSpace(parts[0]),
		Resource: strings.TrimSpace(parts[1]),
		Action:   strings.TrimSpace(parts[2]),
	}, nil
}

// NewEngine loads casbin model+policy and route map.
func NewEngine(cfg EngineConfig) (*Engine, error) {
	if cfg.ModelPath == "" || cfg.PolicyPath == "" || cfg.RoutesPath == "" {
		return nil, errors.New("authz: missing model/policy/routes path")
	}

	m, err := model.NewModelFromFile(cfg.ModelPath)
	if err != nil {
		return nil, errors.Wrap(err, "load model")
	}

	adapter := fileadapter.NewAdapter(cfg.PolicyPath)
	e, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, errors.Wrap(err, "create enforcer")
	}

	// Ensure policies are loaded.
	if err := e.LoadPolicy(); err != nil {
		return nil, errors.Wrap(err, "load policy")
	}

	rules, err := loadRouteRules(cfg.RoutesPath)
	if err != nil {
		return nil, err
	}
	compiled, err := compileRules(rules)
	if err != nil {
		return nil, err
	}

	return &Engine{
		enforcer: e,
		routes:   compiled,
	}, nil
}

// Subject describes caller identity used for authorization.
type Subject struct {
	Role   string
	UserID string
}

var (
	ErrForbidden = errors.New("forbidden")
)

// Authorize checks whether subject may access given request.
// If route is not covered by the rules, the request is allowed (treated as public).
func (e *Engine) Authorize(r *http.Request, sub Subject) error {
	if e == nil || r == nil {
		return nil
	}
	sub.Role = strings.ToLower(strings.TrimSpace(sub.Role))
	rule, params := e.matchRule(strings.ToUpper(r.Method), r.URL.Path)
	if rule == nil {
		return nil // no rule -> public
	}

	perm, err := parsePermission(rule.Permission)
	if err != nil {
		return err
	}

	ownerID := ""
	if rule.Owner != nil {
		switch rule.Owner.Source {
		case ownerSourcePath:
			ownerID = params[rule.Owner.Name]
		case ownerSourceContext:
			ownerID = sub.UserID
		}
	}

	object := perm.Resource
	ok, err := e.enforcer.Enforce(sub.Role, perm.Domain, object, perm.Action, sub.UserID, ownerID)
	if err != nil {
		return errors.Wrap(err, "enforce")
	}
	if !ok {
		return ErrForbidden
	}
	return nil
}

// RequireAuth reports whether the request matches a protected rule (needs authz).
func (e *Engine) RequireAuth(r *http.Request) bool {
	if e == nil || r == nil {
		return false
	}
	rule, _ := e.matchRule(strings.ToUpper(r.Method), r.URL.Path)
	return rule != nil
}

// -------- routes parsing ----------

type OwnerSpec struct {
	Source string `yaml:"source"`
	Name   string `yaml:"name"`
}

type RouteRule struct {
	Method     string     `yaml:"method"`
	Path       string     `yaml:"path"`
	Permission string     `yaml:"permission"`
	Owner      *OwnerSpec `yaml:"owner,omitempty"`
}

type RoutesFile struct {
	Rules []RouteRule `yaml:"rules"`
}

const (
	ownerSourcePath    = "path"
	ownerSourceContext = "context"
)

func loadRouteRules(path string) ([]RouteRule, error) {
	cleanPath := strings.TrimSpace(path)
	if cleanPath == "" {
		return nil, errors.New("empty routes path")
	}
	cleanPath = filepath.Clean(cleanPath)

	// #nosec G304 -- path comes from trusted service config and is normalized above.
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, errors.Wrap(err, "read routes")
	}
	var cfg RoutesFile
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, errors.Wrap(err, "unmarshal routes")
	}
	for i := range cfg.Rules {
		cfg.Rules[i].Method = strings.ToUpper(strings.TrimSpace(cfg.Rules[i].Method))
		cfg.Rules[i].Path = strings.TrimRight(strings.TrimSpace(cfg.Rules[i].Path), "/")
		cfg.Rules[i].Permission = strings.TrimSpace(cfg.Rules[i].Permission)
	}
	return cfg.Rules, nil
}

type compiledRule struct {
	RouteRule
	segments []segment
}

type segment struct {
	Value   string
	IsParam bool
}

func compileRules(rules []RouteRule) ([]compiledRule, error) {
	out := make([]compiledRule, 0, len(rules))
	for _, r := range rules {
		if r.Method == "" || r.Path == "" || r.Permission == "" {
			return nil, fmt.Errorf("invalid rule %+v", r)
		}
		segs := splitPath(r.Path)
		out = append(out, compiledRule{
			RouteRule: r,
			segments:  segs,
		})
	}
	return out, nil
}

func (e *Engine) matchRule(method, path string) (*compiledRule, map[string]string) {
	path = strings.TrimRight(path, "/")
	if path == "" {
		path = "/"
	}
	for i := range e.routes {
		rule := &e.routes[i]
		if rule.Method != method {
			continue
		}
		if params, ok := matchSegments(rule.segments, path); ok {
			return rule, params
		}
	}
	return nil, nil
}

func splitPath(p string) []segment {
	p = strings.TrimPrefix(p, "/")
	if p == "" {
		return nil
	}
	parts := strings.Split(p, "/")
	segs := make([]segment, len(parts))
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			segs[i] = segment{Value: strings.Trim(part, "{}"), IsParam: true}
		} else {
			segs[i] = segment{Value: part}
		}
	}
	return segs
}

func matchSegments(pattern []segment, path string) (map[string]string, bool) {
	path = strings.TrimPrefix(path, "/")
	if path == "" {
		path = "/"
	}
	parts := []string{}
	if path != "/" {
		parts = strings.Split(path, "/")
	}
	if len(pattern) != len(parts) {
		return nil, false
	}
	params := make(map[string]string)
	for i, seg := range pattern {
		val := parts[i]
		if seg.IsParam {
			params[seg.Value] = val
			continue
		}
		if seg.Value != val {
			return nil, false
		}
	}
	return params, true
}
