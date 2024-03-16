package configs

import (
	"filmoteka/pkg/variables"
	"flag"
	"gopkg.in/yaml.v2"
	"os"
)

func ReadAuthAppConfig() (*variables.AuthorizationAppConfig, error) {
	flag.Parse()
	var path string
	flag.StringVar(&path, "auth_config_path", "../../configs/AuthorizationAppConfig.yml", "Путь к конфигу")

	authAppConfig := variables.AuthorizationAppConfig{}
	authAppFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(authAppFile, &authAppConfig)
	if err != nil {
		return nil, err
	}

	return &authAppConfig, nil
}

func ReadRelationalDataBaseConfig() (*variables.RelationalDataBaseConfig, error) {
	flag.Parse()
	var path string
	flag.StringVar(&path, "sql_config_path", "../../configs/AuthorizationSqlDataBaseConfig.yml", "Путь к конфигу")

	relationalDataBaseConfig := variables.RelationalDataBaseConfig{}
	relationalDataBaseFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(relationalDataBaseFile, &relationalDataBaseConfig)
	if err != nil {
		return nil, err
	}

	return &relationalDataBaseConfig, nil
}

func ReadCacheDatabaseConfig() (*variables.CacheDataBaseConfig, error) {
	flag.Parse()
	var path string
	flag.StringVar(&path, "cache_config_path", "../../configs/AuthorizationCacheDataBaseConfig.yml", "Путь к конфигу")

	cacheDataBaseConfig := variables.CacheDataBaseConfig{}
	cacheDataBaseFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(cacheDataBaseFile, &cacheDataBaseConfig)
	if err != nil {
		return nil, err
	}

	return &cacheDataBaseConfig, nil
}
