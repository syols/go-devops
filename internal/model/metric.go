package model

type Metric interface {
	TypeName() string
	String() string
	Payload(name string, key *string) Payload
	FromString(value string) (Metric, error)
	FromPayload(value Payload, key *string) (Metric, error)
}
