package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// type Config struct {
// 	OrgID          string `json:"org_id,omitempty"`
// 	TFS3BucketName string `json:"tfs3_bucket_name,omitempty"`
// 	AccountID      string `json:"account_id,omitempty"`
// 	RoleARN        string `json:"role_arn,omitempty"`
// }

// reason I commented Config struct since I want to be able to create dynamic key-value pairs

type EnvConfigs map[string]map[string]string

func SaveConfig(env string, cfg map[string]string) error {
	var envCfgs EnvConfigs
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return err
		}
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if len(bytes) > 0 {
		if err := json.Unmarshal(bytes, &envCfgs); err != nil {
			return err
		}
	} else {
		envCfgs = make(EnvConfigs)
	}

	envCfgs[env] = cfg

	cfgBytes, err := json.Marshal(envCfgs)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, cfgBytes, 0644)

}

func setEnvVarsFromEnvConfig(env string) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var envCfgs EnvConfigs
	if err := json.Unmarshal(bytes, &envCfgs); err != nil {
		return err
	}
	cfg, ok := envCfgs[env]
	if !ok {
		return fmt.Errorf("no config found for environment: %s", env)
	}

	for key, value := range cfg {
		SetEnv(key, value)
	}
	return nil
}

func SetEnv(key, value string) error {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("export %s=%s", key, value))
	return cmd.Run()
}

func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(homeDir, ".exporter", "config.json")
	return configPath, nil
}

func DeleteConfig(name string) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var envCfgs EnvConfigs
	if err := json.Unmarshal(bytes, &envCfgs); err != nil {
		return err
	}

	_, ok := envCfgs[name]
	if !ok {
		return fmt.Errorf("no config found for environment: %s", name)
	}

	delete(envCfgs, name)
	bytes, err = json.Marshal(envCfgs)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, bytes, 0644)

}

func UpdateConfig(name string, values map[string]string) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var envCfgs EnvConfigs
	if err := json.Unmarshal(bytes, &envCfgs); err != nil {
		return err
	}

	envCfgs[name] = values

	bytes, err = json.Marshal(envCfgs)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, bytes, 0644)
}

func ListEnvironments() error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("")
		return nil
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var envCfgs EnvConfigs
	if err := json.Unmarshal(bytes, &envCfgs); err != nil {
		return err
	}

	for env := range envCfgs {
		fmt.Println(env)
	}
	return nil
}
