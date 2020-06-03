package account

import (
	"configdatabase/src/apimachinery/flowctrl"
	"configdatabase/src/apimachinery/rest"
	"configdatabase/src/apimachinery/util"
	"configdatabase/src/common/auth"
	"configdatabase/src/common/blog"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"sync"
)

const (
	statusSuccess = 0
	authApiPrefix = "/api/v1"
)

const (
	authAppCodeHeaderKey   string = "X-CC-APP-CODE"
	authAppSecretHeaderKey string = "X-CC-APP-SECRET"
	cmdbUser               string = "user"
	cmdbUserID             string = "system"
)

type acDiscovery struct {
	// auth's servers address, must prefixed with http:// or https://
	servers []string
	index   int
	sync.Mutex
}

func (s *acDiscovery) GetServers() ([]string, error) {
	s.Lock()
	defer s.Unlock()

	num := len(s.servers)
	if num == 0 {
		return []string{}, errors.New("oops, there is no auth center server can be used")
	}

	if s.index < num-1 {
		s.index = s.index + 1
		return append(s.servers[s.index-1:], s.servers[:s.index-1]...), nil
	} else {
		s.index = 0
		return append(s.servers[num-1:], s.servers[:num-1]...), nil
	}
}

// NewAccountCenter create a instance to handle resources with AuthCenter.
func NewAccountCenter(tls *util.TLSClientConfig, cfg AuthConfig, reg prometheus.Registerer) (*AccountCenter, error) {
	blog.V(5).Infof("new auth center client with parameters tls: %+v, cfg: %+v", tls, cfg)
	if !auth.IsAuthed() {
		return new(AccountCenter), nil
	}
	client, err := util.NewClient(tls)
	if err != nil {
		return nil, err
	}

	c := &util.Capability{
		Client: client,
		Discover: &acDiscovery{
			servers: cfg.Address,
		},
		Throttle: flowctrl.NewRateLimiter(1000, 1000),
		Mock: util.MockInfo{
			Mocked: false,
		},
		Reg: reg,
	}

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("Accept", "application/json")
	header.Set(authAppCodeHeaderKey, cfg.AppCode)
	header.Set(authAppSecretHeaderKey, cfg.AppSecret)

	return &AccountCenter{
		Config: cfg,
		authClient: &authClient{
			client:      rest.NewRESTClient(c, authApiPrefix),
			Config:      cfg,
			basicHeader: header,
		},
	}, nil
}

type AccountCenter struct {
	Config AuthConfig
	// http client instance
	client rest.ClientInterface
	// http header info
	header     http.Header
	authClient *authClient
}

func (ac *AccountCenter) Enabled() bool {
	return auth.IsAuthed()
}




