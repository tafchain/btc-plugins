package omni

var omniClient *OmniClient

func Init(host, user, pass string, propertyId uint32) {
	omniClient = NewOmniClient(&ConnConfig{Host: host, User: user, Pass: pass}, propertyId)
}

func GetClient() *OmniClient {
	return omniClient
}
