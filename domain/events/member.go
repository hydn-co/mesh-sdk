package events

import (
	"github.com/fgrzl/es"
	"github.com/google/uuid"
)

type MemberJoined struct {
	es.DomainEventBase
	TenantID uuid.UUID `json:"tenant_id"`
	UserID   uuid.UUID `json:"user_id"`
	RoleID   uuid.UUID `json:"role_id"`
}

func (e *MemberJoined) GetDiscriminator() string {
	return "hydn://domain/events/tenants/member_joined"
}

type MemberRoleAssigned struct {
	es.DomainEventBase
	TenantID uuid.UUID `json:"tenant_id"`
	UserID   uuid.UUID `json:"user_id"`
	RoleID   uuid.UUID `json:"role_id"`
}

func (e *MemberRoleAssigned) GetDiscriminator() string {
	return "hydn://domain/events/tenants/member_assigned_role"
}

type MemberRoleUnassigned struct {
	es.DomainEventBase
	TenantID uuid.UUID `json:"tenant_id"`
	UserID   uuid.UUID `json:"user_id"`
	RoleID   uuid.UUID `json:"role_id"`
}

func (e *MemberRoleUnassigned) GetDiscriminator() string {
	return "hydn://domain/events/tenants/member_assigned_role"
}

type MemberLeft struct {
	es.DomainEventBase
	TenantID uuid.UUID `json:"tenant_id"`
	UserID   uuid.UUID `json:"user_id"`
}

func (e *MemberLeft) GetDiscriminator() string {
	return "hydn://domain/events/tenants/member_left"
}
