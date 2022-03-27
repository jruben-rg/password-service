package pwned

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jruben-rg/password-service/go-pwned/config"
)

type (
	PwnedConfig struct {
		Pwned Pwned `yaml:"pwned"`
	}

	Pwned struct {
		Enabled bool          `yaml:"enabled"`
		Timeout time.Duration `yaml:"timeoutSeconds"`
		URL     string        `yaml:"url"`
	}
)

func NewPwnedConfig(filePath string) *Pwned {
	pwnedConfig := &PwnedConfig{}
	err := config.Read(filePath, &pwnedConfig)
	if err != nil {
		panic("Could not read configuration for Pwned endpoint")
	}
	return &pwnedConfig.Pwned
}

func (p Pwned) IsEnabled() bool {
	return p.Enabled
}

func (p Pwned) IsSecurePassword(password string) (bool, error) {

	if p.Enabled {

		// Encode password
		encodedPassword := encodeToSHA1(password)
		passwordPrefix := encodedPassword[:5]
		passwordSuffix := encodedPassword[5:]

		// Create request to pwned service
		pwnedRequest, err := http.NewRequest(http.MethodGet, p.URL+passwordPrefix, nil)
		if err != nil {
			return false, fmt.Errorf("could not create http request to pwned service")
		}

		// Execute request to pwned service
		client := http.Client{}
		responseBody, err := requestPwnedService(&client, pwnedRequest, p.Timeout*time.Second)
		if err != nil {
			return false, fmt.Errorf("%s", err)
		}

		// Find a password suffix match in request results
		found := isPasswordInDictionary(passwordSuffix, responseBody)
		if found {
			return false, nil
		}
	}

	return true, nil
}

func encodeToSHA1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	bytes := h.Sum(nil)

	return strings.ToUpper(fmt.Sprintf("%x", string(bytes)))
}

func requestPwnedService(client *http.Client, request *http.Request, timeout time.Duration) (string, error) {

	timeoutRequest, cancelFunc := context.WithTimeout(request.Context(), timeout)
	defer cancelFunc()
	request = request.WithContext(timeoutRequest)

	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("could not reach pwned service")
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("could not obtain a successful response from pnwed service")
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error when reading pwned response")
	}

	return string(body), nil
}

func isPasswordInDictionary(passwordSuffix string, responseBody string) bool {

	lines := strings.Split(string(responseBody), "\n")
	for _, line := range lines {

		if len(line) > 0 && strings.Contains(line, ":") {
			if result := strings.Compare(passwordSuffix, line[:strings.LastIndex(line, ":")]); result == 0 {
				return true
			}
		}
	}

	return false
}
