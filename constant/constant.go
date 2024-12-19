package constant

var (
	ExchangeEvent = "customer_callbacks"

	// RabbitMQ auth
	UsernameRabbitMQ = "lchh"
	PasswordRabbitMQ = "secret"
	URLRabbitMQ      = "localhost:5671"
	VhostRabbitMQ    = "customers"
	CertPem          = "/home/tls-gen/basic/result/ca_certificate.pem"
	ClientCertPem    = "/home/tls-gen/basic/result/client_DESKTOP-VM3852R_certificate.pem"
	ClientKeyPem     = "/home/tls-gen/basic/result/client_DESKTOP-VM3852R_key.pem"
)
