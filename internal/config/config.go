package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = "/gatorconfig.json"

func Read(configFilePath string) (Config, error) {
	configFile, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, err
	}
	defer configFile.Close()
	buf, err := io.ReadAll(configFile)
	if err != nil {
		return Config{}, err
	}
	var returnVal Config
	if err := json.Unmarshal(buf, &returnVal); err != nil {
		return Config{}, err
	}
	return returnVal, nil
}

func ReadUserConfig() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	returnVal, err := Read(filePath)
	if err != nil {
		return Config{}, err
	}
	return returnVal, nil
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user
	if err := write(c); err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	filePath, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filePath + configFileName, nil
}


func write(c *Config) error {
	cfgFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

    fi, err := os.Lstat(cfgFilePath)

	// Cache old file, in case of malformed config struct?
	// oldContent,err := io.ReadAll(cfgFile)
	// if err != nil {
	//     return err
	// }
	// var cfgOldContents Config
	//
	// if err:= json.Unmarshal(oldContent,&cfgOldContents); err != nil {
	//     return err
	// }

	newContents, err := json.Marshal(c)
	if err != nil {
		return err
	}
    if err = os.WriteFile(cfgFilePath,newContents,fi.Mode().Perm()); err != nil {
        return err
    }

	return nil
}
