package common

import (
	"cloudknox/terraform-provider-cloudknox/utils"
	"errors"
	"os"
	"sync"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	viper "github.com/spf13/viper"
	oauth2 "golang.org/x/oauth2"
)

type ClientParameters struct {
	ServiceAccountId      string
	AccessKey             string
	SecretKey             string
	SharedCredentialsFile string
	Profile               string
}

type Client struct {
	AccessToken string
}

/* Private Variables */
var client *Client
var clientErr error
var clientOnce sync.Once

var config *viper.Viper
var configOnce sync.Once

var logger *log.Logger
var loggerOnce sync.Once

var oAuthConfig *oauth2.Config

const (
	ConfFileDefaultPath = "/conf"
	ConfFileName        = "terraform-provider-cloudknox-config"
)

func getLogger() *log.Logger {
	loggerOnce.Do(
		func() {
			/* Initialize Logger */
			logger = log.New()
			// formatter := &log.TextFormatter{
			// 	DisableColors: false,
			// 	FullTimestamp: true,
			// }
			logger.SetFormatter(&log.JSONFormatter{})
			//ogger.SetLevel(log.DebugLevel)

			err := os.Remove("info.log")

			file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				logger.Fatal(err)
			}

			logger.SetOutput(file)

			logger.Info("Successfully Created Logger Instance!")

			defer file.Close()
		},
	)
	return logger
}

func setConfiguration(parameters *ClientParameters) (*Client, error) {
	configOnce.Do(
		func() {
			/* Initialize Configuration */

			var configurationType string

			log := GetLogger()

			// log.Info(parameters.ServiceAccountId)
			// log.Info(parameters.AccessKey)
			// log.Info(parameters.SecretKey)
			// log.Info(parameters.SharedCredentialsFile)
			// log.Info(parameters.Profile)

			// Configuration Hierarchy
			// Static Credentials
			// Shared Credentials File (Profile)
			// Shared Credentials File (Default)
			// Default Path Credentials File (Profile)
			// Default Path Credentials File (Default)
			// Environment Variables
			// No Credentials Provided Panin

			homedir, _ := homedir.Dir()

			defaultCredentialspath := homedir + "/.cnx/credentials"

			if parameters.Profile == "" {
				parameters.Profile = "Default"
			}

			if parameters.SharedCredentialsFile == "" {
				log.Info("Custom Shared Credentials File not provided")
			} else {
				log.Info("Searching for Shared Credentials File at" + parameters.SharedCredentialsFile)

				if utils.CheckIfPathExists(parameters.SharedCredentialsFile) {
					log.Info("Shared Credentials File located")
					log.Info("Checking " + parameters.Profile + " Profile")
					/*
						If success
						configurationType = "SharedCredentialsFile"
						else
						continue

					*/
				}
			}

			if parameters.ServiceAccountId == "" || parameters.AccessKey == "" || parameters.SecretKey == "" {
				log.Info("Static Credentials not provided or incomplete")
				if configurationType != "" {
					log.Info("Continuing to use " + configurationType)
				}
			} else {

				if configurationType == "" {
					log.Info("Static Credentials provided")

				} else {
					log.Info("Static Credentials provided, overriding " + configurationType)
				}

				configurationType = "StaticCredentials"

				buildClient(parameters.ServiceAccountId, parameters.AccessKey, parameters.SecretKey, configurationType)

				return

			}

			//Check Default Path
			if utils.CheckIfPathExists(defaultCredentialspath) {
				log.Info("Credentials file found at " + defaultCredentialspath)
				log.Info("Checking " + parameters.Profile + " Profile")
				/*
					If success
					configurationType = "DefaultCredentialsFile"
					return
					else
					continue

				*/

				return
			} else {
				log.Info("Checking Environment Variables")

				/*
					If success
					configurationType = "DefaultCredentialsFile"
					return
					else
					error

				*/

				return
			}

			return
		},
	)

	return client, clientErr
}

func buildClient(sai string, ak string, sk string, configurationType string) {
	log := GetLogger()
	log.Info("Using " + configurationType + " to request API Access Token")

	client = &Client{
		AccessToken: "DUMMY VALUE",
	}

	clientErr = errors.New("Client Error")

	return
}

/* Public Variables */

/* Public Functions */
func GetClient() (*Client, error) {
	return client, clientErr
}

func SetConfiguration(parameters *ClientParameters) {
	setConfiguration(parameters)
}

func GetLogger() *logrus.Logger {
	logger := getLogger()

	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logger.Fatal(err)
	}
	logger.SetOutput(file)
	return logger
}