package registrator

//Registrator defines contract for services to register themselves to the Service Discovery Service
type Registrator interface {
	Register(serviceName , serviceID , server string , port int ) error
	Deregister(serviceID string) error
}