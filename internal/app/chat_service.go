package app

import "try_agent_go/internal/domain"

type ChatService struct {
	provider     domain.Provider
	conversation *domain.Conversation
	systemPrompt string
}

func NewChatService(provider domain.Provider, systemPrompt string) *ChatService {
	return &ChatService{
		provider:     provider,
		conversation: domain.NewConversation(),
		systemPrompt: systemPrompt,
	}
}

func (s *ChatService) Send(text string) (string, error) {
	s.conversation.AddUserMessage(text)

	reply, err := s.provider.SendMessage(s.conversation.History(), s.systemPrompt)
	if err != nil {
		s.conversation.UndoLast()
		return "", err
	}

	s.conversation.AddModelMessage(reply)
	return reply, nil
}

func (s *ChatService) SwitchProvider(provider domain.Provider) {
	s.provider = provider
	s.conversation.Reset()
}

func (s *ChatService) ProviderName() string {
	return s.provider.Name()
}
