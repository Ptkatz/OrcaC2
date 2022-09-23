package api

var (
	HOST               = "127.0.0.1:6000"
	CONN_SERVER_API    string
	REGISTER_API       string
	SEND_TO_CLIENT_API string
)

func InitApi(host string) {
	HOST = host
	CONN_SERVER_API = "ws://" + host + "/ws?systemId="
	REGISTER_API = "http://" + host + "/api/register"
	SEND_TO_CLIENT_API = "http://" + host + "/api/send_to_client"
}
