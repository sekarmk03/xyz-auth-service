package authorization

type AccessibleRoles map[string]map[string][]uint32

/*
	1. Super Admin
	2. Admin
	3. Consumer
*/

const (
	BasePath = "xyz-auth-service"
	AuthSvc  = "AuthService"
)

var roles = AccessibleRoles{
	"/" + BasePath + "." + AuthSvc + "/": {
		"GetCurrentUser": {1, 2, 3},
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
