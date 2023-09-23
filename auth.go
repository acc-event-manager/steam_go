package steam_go

// Modified version from github.com/solovev/steam_go
// Credits: github.com/solovev, github.com/ramonberrutti, https://github.com/Fyb3roptik, https://github.com/anotherGoogleFan
import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var (
	steamLogin = "https://steamcommunity.com/openid/login"

	openMode       = "checkid_setup"
	openNS         = "http://specs.openid.net/auth/2.0"
	openIdentifier = "http://specs.openid.net/auth/2.0/identifier_select"

	validationRegexp       = regexp.MustCompile("^(http|https)://steamcommunity.com/openid/id/[0-9]{15,25}$")
	digitsExtractionRegexp = regexp.MustCompile(`\D+`)
)

type OpenID struct {
	root      string
	returnUrl string
	data      url.Values
}

func NewOpenID(r *http.Request) *OpenID {
	id := new(OpenID)
	proto := "http://"
	if r.Header.Get("X-Forwarded-Proto") == "https" || r.TLS != nil {
		proto = "https://"
	}
	if r.Header.Get("X-Forwarded-Host") == "" {
		id.root = proto + r.Host
	} else {
		id.root = proto + r.Header.Get("X-Forwarded-Host")
	}
	uri := r.RequestURI
	if i := strings.Index(uri, "openid."); i != -1 {
		uri = uri[0 : i-1]
	}
	id.returnUrl = id.root + uri
	switch r.Method {
	case "POST":
		r.ParseForm()
		id.data = r.Form
	case "GET":
		id.data = r.URL.Query()
	}
	return id
}

func (id OpenID) AuthUrl() string {
	data := make(url.Values)
	data.Set("openid.claimed_id", openIdentifier)
	data.Set("openid.identity", openIdentifier)
	data.Set("openid.mode", openMode)
	data.Set("openid.ns", openNS)
	data.Set("openid.realm", id.root)
	data.Set("openid.return_to", id.returnUrl)
	url := steamLogin + "?" + data.Encode()
	return url
}

func (id *OpenID) ValidateAndGetID() (string, error) {
	if id.Mode() != "id_res" {
		return "", errors.New("Mode must equal to \"id_res\"")
	}
	// if id.data.Get("openid.return_to") != id.returnUrl {
	// 	return "", errors.New("the \"return_to url\" must match the url of current request")
	// }
	params := make(url.Values)
	params.Set("openid.assoc_handle", id.data.Get("openid.assoc_handle"))
	params.Set("openid.signed", id.data.Get("openid.signed"))
	params.Set("openid.sig", id.data.Get("openid.sig"))
	params.Set("openid.ns", id.data.Get("openid.ns"))
	split := strings.Split(id.data.Get("openid.signed"), ",")
	for _, item := range split {
		params.Set("openid."+item, id.data.Get("openid."+item))
	}
	params.Set("openid.mode", "check_authentication")
	resp, err := http.PostForm(steamLogin, params)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	response := strings.Split(string(content), "\n")
	if response[0] != "ns:"+openNS {
		return "", errors.New("wrong ns in the response")
	}
	if strings.HasSuffix(response[1], "false") {
		return "", errors.New("unable validate openId")
	}
	openIDUrl := id.data.Get("openid.claimed_id")
	if !validationRegexp.MatchString(openIDUrl) {
		return "", errors.New("invalid steam id pattern")
	}
	return digitsExtractionRegexp.ReplaceAllString(openIDUrl, ""), nil
}

func (id OpenID) Mode() string {
	return id.data.Get("openid.mode")
}
