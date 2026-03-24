package main

import (
	"bufio"
	"os"
	"strings"
)

var KBApiKey string
var KBDomain string
var KBServiceID string
var KBSearchServiceID string
var ARKApiKey string
var ArkModelID string
var ServerPort string
var KBAK string
var KBSK string
var KBID string
var AccountID string
var TOSBucket string
var TOSEndpoint string
var TOSRegion string

func LoadConfigFromEnv() {
	loadEnvFile(".env.local")
	loadEnvFile(".env")
	KBApiKey = os.Getenv("KB_API_KEY")
	KBDomain = os.Getenv("KB_DOMAIN")
	if KBDomain == "" {
		KBDomain = "api-knowledgebase.mlp.cn-beijing.volces.com"
	}
	KBServiceID = os.Getenv("KB_SERVICE_ID")
	KBSearchServiceID = os.Getenv("KB_SEARCH_SERVICE_ID")
	ARKApiKey = os.Getenv("ARK_API_KEY")
	ArkModelID = os.Getenv("ARK_MODEL_ID")
	ServerPort = os.Getenv("PORT")
	if ServerPort == "" {
		ServerPort = "8001"
	}
	KBAK = os.Getenv("KB_AK")
	KBSK = os.Getenv("KB_SK")
	KBID = os.Getenv("KB_ID")
	AccountID = os.Getenv("KB_ACCOUNT_ID")
	if AccountID == "" {
		AccountID = os.Getenv("V_ACCOUNT_ID")
	}
	TOSBucket = os.Getenv("TOS_BUCKET")
	TOSEndpoint = os.Getenv("TOS_ENDPOINT")
	if TOSEndpoint == "" {
		TOSEndpoint = "tos-cn-beijing.volces.com"
	}
	TOSRegion = os.Getenv("TOS_REGION")
	if TOSRegion == "" {
		TOSRegion = "cn-beijing"
	}
}

func loadEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		i := strings.Index(line, "=")
		if i <= 0 {
			continue
		}
		k := strings.TrimSpace(line[:i])
		v := strings.TrimSpace(line[i+1:])
		if k != "" {
			os.Setenv(k, v)
		}
	}
}
