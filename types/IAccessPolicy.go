package types

type IAccessPolicy string

const (
	AccessPolicyPublic        IAccessPolicy = "public"
	AccessPolicyAuthenticated IAccessPolicy = "authenticated"
	AccessPolicyAdmin         IAccessPolicy = "admin"
)
