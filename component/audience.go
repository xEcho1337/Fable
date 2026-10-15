package component

// Audience is a destination capable of receiving component messages.
//
// Relay, a server implementation, or a test double owns the concrete
// implementation. Fable only defines the transport-independent contract.
type Audience interface {
	SendMessage(Component) error
}

// SendMessage sends one component to every audience in the collection.
func SendMessage(audiences []Audience, message Component) error {
	for _, audience := range audiences {
		if audience == nil {
			continue
		}

		if err := audience.SendMessage(message); err != nil {
			return err
		}
	}

	return nil
}
