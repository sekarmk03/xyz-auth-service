package authorization

type AccessibleRoles map[string]map[string][]uint32

const (
	BasePath = "xyz-auth-service"
	AuthSvc  = "AuthService"
)

var roles = AccessibleRoles{
	"/" + BasePath + "." + AuthSvc + "/": {
		// "DeletePost":  {1, 2, 8},
	},
}

func GetAccessibleRoles() map[string][]uint32 {
	routes := make(map[string][]uint32)

	for service, methods := range roles {
		for method, methodRoles := range methods {
			route := service + method
			routes[route] = methodRoles
		}
	}

	return routes
}
