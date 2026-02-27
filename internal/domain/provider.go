package domain

type Provider interface {
	SendMessage(history []Message, systemInstruction string) (string, error)
	Name() string
}
