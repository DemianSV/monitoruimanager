package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
)

type (
	// Общие параметры TLS
	TAppConfigTLS struct {
		CERTPATH    string `json:"CERTPATH"`
		KEYPATH     string `json:"KEYPATH"`
		KEYPASSWORD string `json:"KEYPASSWORD"`
		CAPATH      string `json:"CAPATH"`
	}
	// Параметры UIManager
	TAppConfigUIManager struct {
		MODULEID       string   `json:"MODULЕID"`
		SPACEID        string   `json:"SPACEID"`
		DATAMANAGERURL []string `json:"DATAMANAGERURL"`
	}
	// Параметры запуска HTTP сервера
	TAppConfigHTTP struct {
		HOST         string `json:"HOST"`
		PORT         string `json:"PORT"`
		READTIMEOUT  int    `json:"READTIMEOUT"`
		WRITETIMEOUT int    `json:"WRITETIMEOUT"`
		RL           int    `json:"RL"`
		TLS          bool   `json:"TLS"`
		SKIPVERIFY   bool   `json:"SKIPVERIFY"`
		TRUSTPROXY   bool   `json:"TRUSTPROXY"`
	}
	TAppConfigJWT struct {
		SECRET            string `json:"SECRET"`
		ALLOWINSECUREHTTP bool   `json:"ALLOWINSECUREHTTP"`
	}

	/*
		CONSISTENCY:
		Any
		One
		Two
		Three
		Quorum
		All
		LocalQuorum
		EachQuorum
		LocalOn
	*/
	TAppConfigDB struct {
		HOSTS           []string `json:"HOSTS"`
		USERNAME        string   `json:"USERNAME"`
		PASSWORD        string   `json:"PASSWORD"`
		KEYSPACE        string   `json:"KEYSPACE"`
		CONSISTENCY     string   `json:"CONSISTENCY"`
		CONSISTENCYREAD string   `json:"CONSISTENCYREAD"`
		TIMEOUT         int      `json:"TIMEOUT"`
		TLS             bool     `json:"TLS"`
		SKIPVERIFY      bool     `json:"SKIPVERIFY"`
	}

	// Общая структура параметров сервиса
	TAppConfig struct {
		TLS       TAppConfigTLS       `json:"TLS"`
		UIMANAGER TAppConfigUIManager `json:"UIMANAGER"`
		HTTP      TAppConfigHTTP      `json:"HTTP"`
		JWT       TAppConfigJWT       `json:"JWT"`
		DB        TAppConfigDB        `json:"DB"`
	}
)

var (
	AppConfig TAppConfig
)

// LoadConfig Загрузка конфигурации, для Docker версии конфиг собирается из переменного окружения
func (appCfg *TAppConfig) LoadConfig() (err error) {

	log.Print("Загрузка конфигурации из файла...")
	data, err := os.ReadFile(ConfigFile)
	if err == nil {
		err = json.Unmarshal(data, &appCfg)
		if err != nil {
			return err
		}
	} else {
		log.Print("Не удалось прочитать конфигурацию из файла (" + err.Error() + ")!")
	}

	/*
		Переменное окружение (для Docker версии):
		MONITOR_TLS_CERTPATH Путь к сертификату для TLS взаимодействия
		MONITOR_TLS_KEYPATH Путь к приватному ключу для TLS взаимодействия
		MONITOR_TLS_KEYPASSWORD Пароль для приватного ключа TLS
		MONITOR_TLS_CAPATH Путь к CA сертификату для TLS взаимодействия

		MONITOR_UIMANAGER_MODULEID ID Модуля зарезервированный для UIManager
		MONITOR_UIMANAGER_SPACEID ID Пространства зарезервированного для UIManager
		MONITOR_UIMANAGER_DATAMANAGERURL URL API DataManager

		MONITOR_HTTP_HOST Адрес для входящих соединений API
		MONITOR_HTTP_PORT Порт для входящих соединений API
		MONITOR_HTTP_READTIMEOUT Read Timeout для HTTP сервера
		MONITOR_HTTP_WRITETIMEOUT Write Timeout для HTTP сервера
		MONITOR_HTTP_RL Requests Limit
		MONITOR_HTTP_TLS Включение TLS взаимодействия с http сервером (true/false)
		MONITOR_HTTP_SKIPVERIFY Выключение проверки сертификата (true/false)
		MONITOR_HTTP_TRUSTPROXY Доверять X-Forwarded-Proto от reverse proxy (true/false)

		MONITOR_JWT_SECRET Секрет для подписи токенов (минимум 32 байта)
		MONITOR_JWT_ALLOW_INSECURE_HTTP Разрешить передачу пароля по HTTP (только тестовые окружения)

		MONITOR_DB_HOST Массив адресов для подключения к БД, разделитель запятая (127.0.0.1:19042,127.0.0.1:29042)
		MONITOR_DB_USERNAME Имя пользователя для подключения к БД
		MONITOR_DB_PASSWORD Пароль для подключения к БД
		MONITOR_DB_TLS Включение TLS взаимодействия с кластером БД (true/false)
		MONITOR_DB_HOSTVERIFICATION Включение проверки сертификата при подключении к БД (true/false)
		MONITOR_DB_KEYSPACE Пространство ключей (KeySpace) для подключения
		MONITOR_DB_CONSISTENCY Параметр Consistency для данных в кластере БД
		MONITOR_DB_CONSISTENCYREAD Параметр Consistency для выборки данных из БД
		MONITOR_DB_TIMEOUT Значение времени ожидания выполнения запроса до принудительного завершения
	*/

	// TLS общие параметры
	log.Print("Загрузка конфигурации из переменного окружения...")
	var tlsCertPathENV string
	tlsCertPathENV, _ = os.LookupEnv("MONITOR_TLS_CERTPATH")
	if tlsCertPathENV != "" {
		appCfg.TLS.CERTPATH = tlsCertPathENV
	}
	var tlsKeyPathENV string
	tlsKeyPathENV, _ = os.LookupEnv("MONITOR_TLS_KEYPATH")
	if tlsKeyPathENV != "" {
		appCfg.TLS.KEYPATH = tlsKeyPathENV
	}
	var tlsKeyPasswordENV string
	tlsKeyPasswordENV, _ = os.LookupEnv("MONITOR_TLS_KEYPASSWORD")
	if tlsKeyPasswordENV != "" {
		appCfg.TLS.KEYPASSWORD = tlsKeyPasswordENV
	}
	var tlsCAPathENV string
	tlsCAPathENV, _ = os.LookupEnv("MONITOR_TLS_CAPATH")
	if tlsCAPathENV != "" {
		appCfg.TLS.CAPATH = tlsCAPathENV
	}

	// UIManager
	var uimanagerModuleIDENV string
	uimanagerModuleIDENV, _ = os.LookupEnv("MONITOR_UIMANAGER_MODULEID")
	if uimanagerModuleIDENV != "" {
		appCfg.UIMANAGER.MODULEID = uimanagerModuleIDENV
	}
	var uimanagerSpaceIDENV string
	uimanagerSpaceIDENV, _ = os.LookupEnv("MONITOR_UIMANAGER_SPACEID")
	if uimanagerSpaceIDENV != "" {
		appCfg.UIMANAGER.SPACEID = uimanagerSpaceIDENV
	}
	var uimanagerDataManagerURLENV string
	uimanagerDataManagerURLENV, _ = os.LookupEnv("MONITOR_UIMANAGER_DATAMANAGERURL")
	if uimanagerDataManagerURLENV != "" {
		uimanagerDataManagerURLENV = strings.TrimSpace(uimanagerDataManagerURLENV)       // Удаление возможных пробелов
		appCfg.UIMANAGER.DATAMANAGERURL = strings.Split(uimanagerDataManagerURLENV, ",") // Разбиваем строку по символу ","
	}

	// HTTP сервер
	var httpHostENV string
	httpHostENV, _ = os.LookupEnv("MONITOR_HTTP_HOST")
	if httpHostENV != "" {
		appCfg.HTTP.HOST = httpHostENV
	}
	var httpPortENV string
	httpPortENV, _ = os.LookupEnv("MONITOR_HTTP_PORT")
	if httpPortENV != "" {
		appCfg.HTTP.PORT = httpPortENV
	}
	var httpReadTimeoutENV string
	httpReadTimeoutENV, _ = os.LookupEnv("MONITOR_HTTP_READTIMEOUT")
	if httpReadTimeoutENV != "" {
		var httpReadTimeoutINT int
		httpReadTimeoutINT, err = strconv.Atoi(httpReadTimeoutENV)
		if err != nil {
			log.Print("Не удалось конвертировать MONITOR_HTTP_READTIMEOUT (" + err.Error() + ")!")
		} else {
			appCfg.HTTP.READTIMEOUT = httpReadTimeoutINT
		}
	}
	var httpWriteTimeoutENV string
	httpWriteTimeoutENV, _ = os.LookupEnv("MONITOR_HTTP_WRITETIMEOUT")
	if httpWriteTimeoutENV != "" {
		var httpWriteTimeoutINT int
		httpWriteTimeoutINT, err = strconv.Atoi(httpWriteTimeoutENV)
		if err != nil {
			log.Print("Не удалось конвертировать MONITOR_HTTP_WRITETIMEOUT (" + err.Error() + ")!")
		} else {
			appCfg.HTTP.WRITETIMEOUT = httpWriteTimeoutINT
		}
	}
	var httpRLENV string
	httpRLENV, _ = os.LookupEnv("MONITOR_HTTP_RL")
	if httpRLENV != "" {
		var httpRLINT int
		httpRLINT, err = strconv.Atoi(httpRLENV)
		if err != nil {
			log.Print("Не удалось конвертировать MONITOR_HTTP_RL (" + err.Error() + ")!")
		} else {
			appCfg.HTTP.RL = httpRLINT
		}
	}
	var httpTLSENV string
	httpTLSENV, _ = os.LookupEnv("MONITOR_HTTP_TLS")
	if httpTLSENV != "" {
		boolValue, err := strconv.ParseBool(httpTLSENV)
		if err != nil {
			log.Print("Не удалось конвертировать MONITOR_HTTP_TLS (" + err.Error() + ")!")
		} else {
			appCfg.HTTP.TLS = boolValue
		}
	}
	httpTrustProxyENV, _ := os.LookupEnv("MONITOR_HTTP_TRUSTPROXY")
	if httpTrustProxyENV != "" {
		boolValue, err := strconv.ParseBool(httpTrustProxyENV)
		if err != nil {
			log.Print("Не удалось конвертировать MONITOR_HTTP_TRUSTPROXY (", err.Error(), ")!")
		} else {
			appCfg.HTTP.TRUSTPROXY = boolValue
		}
	}
	var httpSkipVerify string
	httpSkipVerify, _ = os.LookupEnv("MONITOR_HTTP_TLS_SKIPVERIFY")
	if httpSkipVerify != "" {
		boolValue, err := strconv.ParseBool(httpSkipVerify)
		if err != nil {
			log.Print("Не удалось конвертировать MONITOR_HTTP_TLS_SKIPVERIFY (" + err.Error() + ")!")
		} else {
			appCfg.HTTP.SKIPVERIFY = boolValue
		}
	}
	jwtSecretENV, _ := os.LookupEnv("MONITOR_JWT_SECRET")
	if jwtSecretENV != "" {
		appCfg.JWT.SECRET = jwtSecretENV
	}
	jwtAllowInsecureHTTPENV, _ := os.LookupEnv("MONITOR_JWT_ALLOW_INSECURE_HTTP")
	if jwtAllowInsecureHTTPENV != "" {
		boolValue, err := strconv.ParseBool(jwtAllowInsecureHTTPENV)
		if err != nil {
			log.Print("Не удалось конвертировать MONITOR_JWT_ALLOW_INSECURE_HTTP (", err.Error(), ")!")
		} else {
			appCfg.JWT.ALLOWINSECUREHTTP = boolValue
		}
	}

	// Подключение к БД
	var dbHostsENV string
	dbHostsENV, _ = os.LookupEnv("MONITOR_DB_HOST")
	if dbHostsENV != "" {
		dbHostsENV = strings.TrimSpace(dbHostsENV)       // Удаление возможных пробелов
		appCfg.DB.HOSTS = strings.Split(dbHostsENV, ",") // Разбиваем строку по символу ","
	}
	var dbUserNameENV string
	dbUserNameENV, _ = os.LookupEnv("MONITOR_DB_USERNAME")
	if dbUserNameENV != "" {
		appCfg.DB.USERNAME = dbUserNameENV
	}
	var dbPasswordENV string
	dbPasswordENV, _ = os.LookupEnv("MONITOR_DB_PASSWORD")
	if dbPasswordENV != "" {
		appCfg.DB.PASSWORD = dbPasswordENV
	}
	var dbKeySpaceENV string
	dbKeySpaceENV, _ = os.LookupEnv("MONITOR_DB_KEYSPACE")
	if dbKeySpaceENV != "" {
		appCfg.DB.KEYSPACE = dbKeySpaceENV
	}
	var dbConsistensyENV string
	dbConsistensyENV, _ = os.LookupEnv("MONITOR_DB_CONSISTENCY")
	if dbConsistensyENV != "" {
		appCfg.DB.CONSISTENCY = dbConsistensyENV
	}
	var dbConsistensyReadENV string
	dbConsistensyReadENV, _ = os.LookupEnv("MONITOR_DB_CONSISTENCYREAD")
	if dbConsistensyReadENV != "" {
		appCfg.DB.CONSISTENCYREAD = dbConsistensyReadENV
	}
	var dbTimeoutENV string
	dbTimeoutENV, _ = os.LookupEnv("MONITOR_DB_TIMEOUT")
	if dbTimeoutENV != "" {
		var dbTimeoutINT int
		dbTimeoutINT, err = strconv.Atoi(dbTimeoutENV)
		if err != nil {
			log.Print("Не удалось конвертировать MONITOR_DB_TIMEOUT (" + err.Error() + ")!")
		} else {
			appCfg.DB.TIMEOUT = dbTimeoutINT
		}
	}
	var dbTLSENV string
	dbTLSENV, _ = os.LookupEnv("MONITOR_DB_TLS")
	if dbTLSENV != "" {
		boolValue, err := strconv.ParseBool(dbTLSENV)
		if err != nil {
			log.Print("Не удалось конвертировать MONITOR_DB_TLS (" + err.Error() + ")!")
		} else {
			appCfg.DB.TLS = boolValue
		}
	}
	var dbSkipVerifyENV string
	dbSkipVerifyENV, _ = os.LookupEnv("MONITOR_DB_TLS_SKIPVERIFY")
	if dbSkipVerifyENV != "" {
		boolValue, err := strconv.ParseBool(dbSkipVerifyENV)
		if err != nil {
			log.Print("Не удалось конвертировать MONITOR_DB_TLS_SKIPVERIFY (" + err.Error() + ")!")
		} else {
			appCfg.DB.SKIPVERIFY = boolValue
		}
	}

	log.Print("Конфигурационные параметры:")

	log.Print("TLS.CERTPATH: ", appCfg.TLS.CERTPATH)
	log.Print("TLS.KEYPATH: ", appCfg.TLS.KEYPATH)
	log.Print("TLS.CAPATH: ", appCfg.TLS.CAPATH)

	log.Print("UIMANAGER.MODULEID: ", appCfg.UIMANAGER.MODULEID)
	log.Print("UIMANAGER.SPACEID: ", appCfg.UIMANAGER.SPACEID)
	log.Print("UIMANAGER.DATAMANAGERURL: ", appCfg.UIMANAGER.DATAMANAGERURL)

	log.Print("HTTP.HOST: ", appCfg.HTTP.HOST)
	log.Print("HTTP.PORT: ", appCfg.HTTP.PORT)
	log.Print("HTTP.READTIMEOUT: ", appCfg.HTTP.READTIMEOUT)
	log.Print("HTTP.WRITETIMEOUT: ", appCfg.HTTP.WRITETIMEOUT)
	log.Print("HTTP.RL: ", appCfg.HTTP.RL)
	log.Print("HTTP.TLS: ", appCfg.HTTP.TLS)
	log.Print("HTTP.SKIPVERIFY: ", appCfg.HTTP.SKIPVERIFY)
	log.Print("HTTP.TRUSTPROXY: ", appCfg.HTTP.TRUSTPROXY)
	log.Print("JWT.SECRET configured: ", appCfg.JWT.SECRET != "")
	log.Print("JWT.ALLOWINSECUREHTTP: ", appCfg.JWT.ALLOWINSECUREHTTP)

	log.Print("DB.HOSTS: ", appCfg.DB.HOSTS)
	log.Print("DB.USERNAME: ", appCfg.DB.USERNAME)
	log.Print("DB.TLS: ", appCfg.DB.TLS)
	log.Print("DB.SKIPVERIFY: ", appCfg.DB.SKIPVERIFY)
	log.Print("DB.KEYSPACE: ", appCfg.DB.KEYSPACE)
	log.Print("DB.CONSISTENCY: ", appCfg.DB.CONSISTENCY)
	log.Print("DB.CONSISTENCYREAD: ", appCfg.DB.CONSISTENCYREAD)
	log.Print("DB.TIMEOUT: ", appCfg.DB.TIMEOUT)

	return nil
}
