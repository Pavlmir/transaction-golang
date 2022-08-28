package get_settings

import (
    "fmt"
    "os"
    "path/filepath"
	"encoding/json"
)

type Settings struct {
	DBHost     string
	DBPort     string
	DBUsername string
	DBPassword string
	DBName     string
}

func GetFile(path string) *os.File {
	ex, err := os.Executable()
	if err != nil {
		fmt.Printf("Ошибка определения текущего пути %s", err)
		panic(err)
	}
	exPath := filepath.Dir(ex)

	file, err := os.Open(exPath + path + "/common/config.json")
	if err != nil {
		fmt.Printf("Ошибка получения файла %s", err)
		os.Exit(1)
	}

	return file
}

func GetSettings(path string) Settings {
	var settings Settings
      
	jsonParser := json.NewDecoder(GetFile(path))
	if err := jsonParser.Decode(&settings); err != nil {
		fmt.Printf("Не удалось загрузить конфигурационный файл %s", err.Error())
	}
    
    return settings
}
