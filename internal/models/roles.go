package models

// User roles
const (
	RoleAdmin    = "admin"
	RoleGymOwner = "gym_owner"
	RoleTrainer  = "trainer"
	RoleCustomer = "customer"
)

// Role represents a user role in the system
type Role struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

// Define permissions
const (
	PermissionCreateTrainer  = "create:trainer"
	PermissionUpdateTrainer  = "update:trainer"
	PermissionDeleteTrainer  = "delete:trainer"
	PermissionViewTrainer    = "view:trainer"
	PermissionCreateCustomer = "create:customer"
	PermissionUpdateCustomer = "update:customer"
	PermissionDeleteCustomer = "delete:customer"
	PermissionViewCustomer   = "view:customer"
)

// RolePermissions maps roles to their permissions
var RolePermissions = map[string][]string{
	RoleAdmin: {
		PermissionCreateTrainer, PermissionUpdateTrainer, PermissionDeleteTrainer, PermissionViewTrainer,
		PermissionCreateCustomer, PermissionUpdateCustomer, PermissionDeleteCustomer, PermissionViewCustomer,
	},
	RoleGymOwner: {
		PermissionCreateTrainer, PermissionUpdateTrainer, PermissionDeleteTrainer, PermissionViewTrainer,
		PermissionCreateCustomer, PermissionUpdateCustomer, PermissionDeleteCustomer, PermissionViewCustomer,
	},
	RoleTrainer: {
		PermissionViewCustomer,
	},
	RoleCustomer: {},
}
