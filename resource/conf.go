package resource

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var ExePath string
var Path = "resource/static/config.yml"
var Conf = new(ProfileInfo)

type ProfileInfo struct {
	*App              `mapstructure:"app"`
	*RedisConfig      `mapstructure:"redis"`
	*MysqlConfig      `mapstructure:"mysql"`
	*JwtConfig        `mapstructure:"jwt"`
	*DllConfig        `mapstructure:"dll"`
	*ObjStorageConfig `mapstructure:"obj-storage"`
}

// 系统配置
type App struct {
	Env         string          `mapstructure:"env" yaml:"env"`
	Cache       bool            `mapstructure:"cache" yaml:"cache"`
	ServiceName string          `mapstructure:"service-name" yaml:"service-name"`
	MachineID   int64           `mapstructure:"machine-id" yaml:"machine-id"`
	ServerPort  int             `mapstructure:"server-port" yaml:"server-port"`
	ApiPrefix   string          `mapstructure:"api-prefix" yaml:"api-prefix"`
	PrivateKey  string          `mapstructure:"private-key" yaml:"private-key"` // RSA私钥文件路径
	privateKey  *rsa.PrivateKey // 解析后的RSA私钥
}

// GetPrivateKey 获取RSA私钥
func (a *App) GetPrivateKey() *rsa.PrivateKey {
	return a.privateKey
}

// SetPrivateKey 设置RSA私钥
func (a *App) SetPrivateKey(key *rsa.PrivateKey) {
	a.privateKey = key
}

// Redis配置
type RedisConfig struct {
	Addr     string `mapstructure:"addr" yaml:"addr"`
	Password string `mapstructure:"password" yaml:"password"`
	DB       int    `mapstructure:"db" yaml:"db"`
}

// mysql配置
type MysqlConfig struct {
	Host         string `mapstructure:"host" yaml:"host"`
	Port         int    `mapstructure:"port" yaml:"port"`
	User         string `mapstructure:"user" yaml:"user"`
	Password     string `mapstructure:"password" yaml:"password"`
	DBName       string `mapstructure:"dbname" yaml:"dbname"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" yaml:"max-idle-conns"`
	MaxOpenConns int    `mapstructure:"max-open-conns" yaml:"max-open-conns"`
}

// jwt配置
type JwtConfig struct {
	AccessExpire       int64  `mapstructure:"access-expire" yaml:"access-expire"`
	RefreshExpire      int64  `mapstructure:"refresh-expire" yaml:"refresh-expire"`
	Issuer             string `mapstructure:"issuer" yaml:"issuer"`
	AccessTokenSecret  string `mapstructure:"asecret" yaml:"asecret"`
	RefreshTokenSecret string `mapstructure:"rsecret" yaml:"rsecret"`
}

// dll加载配置
type DllConfig struct {
	DllPath string `mapstructure:"dll-path" yaml:"dll-path"`
	DllName string `mapstructure:"dll-name" yaml:"dll-name"`
}

type ObjStorageConfig struct {
	Region       string `mapstructure:"region" yaml:"region"`
	BucketName   string `mapstructure:"bucket-name" yaml:"bucket-name"`
	AccessID     string `mapstructure:"access-id" yaml:"access-id"`
	AccessSecret string `mapstructure:"access-secret" yaml:"access-secret"`
	RoleArn      string `mapstructure:"role-arn" yaml:"role-arn"`
}

// ConfigInit 初始化配置
// 将配置文件的信息反序列化到结构体中
func ConfigInit() {
	// 获取当前工作目录作为项目根目录
	workDir, e := os.Getwd()
	if e != nil {
		fmt.Println("Error:", e)
		return
	}
	ExePath = workDir
	Path = filepath.Join(ExePath, Path)

	configFile := Path
	s := flag.String("f", configFile, "choose conf file.")
	flag.Parse()
	//viper.AddConfigPath(configPath)
	//viper.SetConfigName("conf")     // 读取配置文件
	viper.SetConfigFile(*s)     // 读取配置文件
	err := viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig() faild error:%v\n", err)
		return
	}
	// 把读取到的信息反序列化到Conf变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed,err:%v\n", err)
		return
	}

	// 读取并解析私钥文件
	if err := loadPrivateKey(); err != nil {
		log.Fatalf("load private key failed: %v\n", err)
		return
	}

	viper.WatchConfig() // （热加载时读取配置）监控配置文件
	viper.OnConfigChange(func(in fsnotify.Event) { // 配置文件修改时触发回调
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed,err:%v\n", err)
			return
		}
		// 重新加载私钥
		if err := loadPrivateKey(); err != nil {
			fmt.Printf("reload private key failed: %v\n", err)
		}
	})
}

// loadPrivateKey 加载RSA私钥
func loadPrivateKey() error {
	// 从配置中获取私钥文件路径
	privateKeyPath := filepath.Join(ExePath, Conf.App.PrivateKey)

	// 读取私钥文件
	privateKeyPEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read private key file: %v", err)
	}

	// 解析私钥
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return fmt.Errorf("failed to parse private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}

	// 保存解析后的私钥
	Conf.App.SetPrivateKey(privateKey)
	return nil
}

// SaveConfig 将配置写入文件，并在失败时回滚
func SaveConfig(configFile string, config *ProfileInfo) error {
	// 备份当前配置文件
	backupFile := configFile + ".bak"
	if err := backupConfigFile(configFile, backupFile); err != nil {
		return fmt.Errorf("failed to backup conf file: %v", err)
	}

	// 将 ProfileInfo 转换为 map(识别mapstructure标签)
	configMap := make(map[string]interface{})
	err := mapstructure.Decode(config, &configMap)
	if err != nil {
		rollbackConfigFile(configFile, backupFile)
		return fmt.Errorf("failed to decode conf to map: %v", err)
	}

	// 生成 YAML 数据
	data, err := yaml.Marshal(configMap)
	if err != nil {
		rollbackConfigFile(configFile, backupFile)
		return fmt.Errorf("failed to marshal conf data: %v", err)
	}

	// 写入配置文件
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		rollbackConfigFile(configFile, backupFile)
		return fmt.Errorf("failed to write conf file: %v", err)
	}

	return nil
}

// 备份配置文件
func backupConfigFile(configFile, backupFile string) error {
	// 读取原始配置文件内容
	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果配置文件不存在，不需要备份，直接返回 nil
			return nil
		}
		return err
	}

	// 写入备份文件
	if err := os.WriteFile(backupFile, data, 0644); err != nil {
		return err
	}

	return nil
}

// 回滚配置文件
func rollbackConfigFile(configFile, backupFile string) {
	// 恢复备份文件内容到原始配置文件
	data, err := os.ReadFile(backupFile)
	if err != nil {
		fmt.Printf("error: failed to read backup file %s: %v\n", backupFile, err)
		return
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		fmt.Printf("error: failed to restore conf file from backup %s: %v\n", backupFile, err)
	}

	// 删除备份文件
	if err := os.Remove(backupFile); err != nil {
		fmt.Printf("warning: failed to remove backup file %s: %v\n", backupFile, err)
	}
}
