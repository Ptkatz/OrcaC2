package api

var (
	HOST               string
	CONN_SERVER_API    string
	REGISTER_API       string
	SEND_TO_CLIENT_API string
	MASTER_LOGIN_API   string
)

func InitApi(host string) {
	HOST = host
	CONN_SERVER_API = "ws://" + host + "/ws?systemId="
	REGISTER_API = "http://" + host + "/api/register"
	SEND_TO_CLIENT_API = "http://" + host + "/api/send_to_client"
	MASTER_LOGIN_API = "http://" + host + "/api/master_login"
}
