package securityConfig

import (
	roleAuth "JwtAuth/src/authorization/builder/roleBuilder"
)

type WebSecurityConfig struct {
}

func (s WebSecurityConfig) Config(security *roleAuth.RoleConfigurer) {
	security.
		AntMatcher("/api/v1/test").HasAnyRoles("Admin").
		AntMatcher("/api/v1/**").DenyAll().
		AnyRequest().DenyAll()
}