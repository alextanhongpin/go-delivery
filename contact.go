package omni

// ContactMechanisms represents the identity of the contact, used to represent
// the sender/receiver or from/to pair.
type ContactMechanism struct {
	Name  string // John Doe, or email if not present.
	Type  string // email | mobile
	Value string // john.doe@mail.com
}
