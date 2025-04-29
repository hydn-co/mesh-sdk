package events

import (
	"github.com/fgrzl/es"
	"github.com/hydn-co/mesh-sdk/domain/permissions"
)

type RoleAdded struct {
	es.DomainEventBase
	RoleName    string                   `json:"role_name"`
	Permissions []permissions.Permission `json:"permissions"`
}

func (e *RoleAdded) GetDiscriminator() string {
	return "hydn://domain/events/tenants/role_defined"
}

type RoleRemoved struct {
	es.DomainEventBase
	RoleName    string                   `json:"role_name"`
	Permissions []permissions.Permission `json:"permissions"`
}

func (e *RoleRemoved) GetDiscriminator() string {
	return "hydn://domain/events/tenants/role_defined"
}

type RolePermissionAdded struct {
	es.DomainEventBase
	Permissions permissions.Permission `json:"permission"`
}

func (e *RolePermissionAdded) GetDiscriminator() string {
	return "hydn://domain/events/tenants/role_permissions_changed"
}

type RolePermissionRemoved struct {
	es.DomainEventBase
	Permissions permissions.Permission `json:"permission"`
}

func (e *RolePermissionRemoved) GetDiscriminator() string {
	return "hydn://domain/events/tenants/role_permissions_changed"
}
