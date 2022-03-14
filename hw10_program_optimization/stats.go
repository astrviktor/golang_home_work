package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domain = "." + domain
	scanner := bufio.NewScanner(r)

	result := make(DomainStat)

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	var user User
	for scanner.Scan() {
		bytes := scanner.Bytes()

		if strings.Contains(string(bytes), domain) {
			user = User{}
			if err := json.Unmarshal(bytes, &user); err != nil {
				return DomainStat{}, err
			}

			matched := strings.HasSuffix(user.Email, domain)

			if matched {
				num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
				num++
				result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
			}
		}
	}

	return result, nil
}

func GetDomainStatSlow(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		var user User
		if err = json.Unmarshal([]byte(line), &user); err != nil {
			return
		}
		result[i] = user
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched, err := regexp.Match("\\."+domain, []byte(user.Email))
		if err != nil {
			return nil, err
		}

		if matched {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}
