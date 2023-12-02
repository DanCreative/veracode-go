package veracode

import (
	"errors"
	"fmt"
	"os"
	"path"

	"gopkg.in/ini.v1"
)

// getCredentialsFilePath gets the Veracode API credentials file path.
func getCredentialsFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(homeDir, ".veracode", "credentials"), nil
}

// GetProfile returns a pointer to the section of the credentials file that the user is using.
func getProfile(filePath string) (*ini.Section, error) {
	profile := os.Getenv("VERACODE_API_PROFILE")

	cfg, err := ini.Load(filePath)
	if err != nil {
		return nil, fmt.Errorf("error loading ini file. Message: %s", err.Error())
	}

	var rSection *ini.Section
	sections := cfg.Sections()

	if len(sections) == 1 && sections[0].Name() == "DEFAULT" {
		rSection = sections[0]
	} else {
		// If profile is set, try and get it from ini file.
		// If profile is not set, try and get the default from ini file.
		if profile != "" {
			rSection, err = cfg.GetSection(profile)
			if err != nil {
				return nil, fmt.Errorf("error loading profile: %s from file. Message: %s", profile, err.Error())
			}
		} else {
			rSection, err = cfg.GetSection("default")
			if err != nil {
				return nil, fmt.Errorf("no profile set and error loading default profile. Message: %s", err.Error())
			}
		}
	}
	return rSection, err
}

// LoadVeracodeCredentails will get the Veracode API key and secret for set profile from the credentials file.
// The profile name will be read from the VERACODE_API_PROFILE environmental variable. If the variable is not set, the
// profile with name "default" will be used. If there is only one profile with no name it will be used.
// The credentials file should be in the .ini format and should be present in the /.veracode/ folder in the user's home
// directory. Please refer to the documentation for more information: https://docs.veracode.com/r/c_httpie_tool.
func LoadVeracodeCredentials() (string, string, error) {
	credsPath, err := getCredentialsFilePath()
	if err != nil {
		return "", "", err
	}

	profile, err := getProfile(credsPath)
	if err != nil {
		return "", "", err
	}

	key, secret := profile.Key("veracode_api_key_id").String(), profile.Key("veracode_api_key_secret").String()
	if key == "" || secret == "" {
		err := errors.New("failed to load Veracode API credentials from file. Please refer to documentation: https://docs.veracode.com/r/c_httpie_tool")
		return "", "", err
	}

	return key, secret, nil
}
