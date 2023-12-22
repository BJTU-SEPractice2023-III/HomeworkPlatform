package bootstrap

import (
	"homework_platform/internal/config"
	"os"
)

var Config *config.Config

func InitConfig() {
	var err error
	// log.Println("[bootstrap/InitConfig]: Initializing config...")
	Config = config.DefaultACHConfig()

	Config, err = config.ReadConfig()
	if err != nil {
		if os.IsNotExist(err) { // 文件不存在，创建并写入默认配置
			println("[bootstrap/InitConfig]: Cannot find config.yml, creating default config...")
			config.Save(Config)
			println("[bootstrap/InitConfig]: Successfully created config.yml, please complete the config.")
		}
		os.Exit(1)
	}

	// // log.Print("[bootstrap/InitConfig]: config: ", Config, '\n')
	// os.Mkdir(Config.BackupDir, 0666)
	// os.Mkdir(path.Join(Config.BackupDir, "backups"), 0666)
	// os.Mkdir(path.Join(Config.BackupDir, "snapshots"), 0666)
}
