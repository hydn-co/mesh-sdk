package entities

import (
	"github.com/fgrzl/json/polymorphic"
	"github.com/google/uuid"
)

func init() {
	polymorphic.Register(func() *NetworkObject { return &NetworkObject{} })
}

func NewNetworkObject() *NetworkObject { return &NetworkObject{} }

type IPType string

const (
	NetworkObjectDiscriminator = "hydn://network_objects"
	NetworkObjectSpace         = "network_objects"

	IPTypeUndefined       IPType = "UNDEFINED"
	IPTypeV4Unicast       IPType = "V4_UNICAST"
	IPTypeV4Broadcast     IPType = "V4_BROADCAST"
	IPTypeV4Multicast     IPType = "V4_MULTICAST"
	IPTypeV4Private       IPType = "V4_PRIVATE"
	IPTypeV4Public        IPType = "V4_PUBLIC"
	IPTypeV4Loopback      IPType = "V4_LOOPBACK"
	IPTypeV6Unicast       IPType = "V6_UNICAST"
	IPTypeV6Multicast     IPType = "V6_MULTICAST"
	IPTypeV6Anycast       IPType = "V6_ANYCAST"
	IPTypeV6LinkLocal     IPType = "V6_LINK_LOCAL"
	IPTypeV6GlobalUnicast IPType = "V6_GLOBAL_UNICAST"
	IPTypeV6UniqueLocal   IPType = "V6_UNIQUE_LOCAL"
	IPTypeV6Loopback      IPType = "V6_LOOPBACK"
)

type OperationalStatus string

const (
	OpStatusUnknown  OperationalStatus = "UNKNOWN"
	OpStatusActive   OperationalStatus = "ACTIVE"
	OpStatusInactive OperationalStatus = "INACTIVE"
	OpStatusFaulty   OperationalStatus = "FAULTY"
)

type IPAddress struct {
	Address string `json:"address"`
	IPType  IPType `json:"type" enums:"UNDEFINED,V4_UNICAST,V4_BROADCAST,V4_MULTICAST,V4_PRIVATE,V4_PUBLIC,V4_LOOPBACK,V6_UNICAST,V6_MULTICAST,V6_ANYCAST,V6_LINK_LOCAL,V6_GLOBAL_UNICAST,V6_UNIQUE_LOCAL,V6_LOOPBACK"`
}

type NetworkObject struct {
	ConnectorID            uuid.UUID         `json:"connector_id"`
	NetworkObjectReference string            `json:"network_object_reference"`
	Type                   string            `json:"type"`
	Name                   string            `json:"name"`
	IPAddresses            []IPAddress       `json:"ip_addresses"`
	MACAddress             string            `json:"mac_address,omitempty"`
	Tags                   []string          `json:"tags,omitempty"`
	Location               string            `json:"location,omitempty"`
	Status                 OperationalStatus `json:"status" enums:"UNKNOWN,ACTIVE,INACTIVE,FAULTY"`
	Metadata               map[string]string `json:"metadata,omitempty"`
}

func (e *NetworkObject) GetConnectorID() uuid.UUID {
	return e.ConnectorID
}

func (e *NetworkObject) GetDistinctID() uuid.UUID {
	return uuid.NewSHA1(e.ConnectorID, []byte(e.GetReference()))
}

func (e *NetworkObject) GetReference() string {
	return e.NetworkObjectReference
}

func (e *NetworkObject) GetSpace() string {
	return "network-objects"
}

func (e *NetworkObject) GetDiscriminator() string {
	return "hydn://catalog/v1/entities/network-object"
}
