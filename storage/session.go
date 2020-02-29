package storage

import (
	"animated-robot/domain"
	"animated-robot/tools"
	"fmt"
	"github.com/google/uuid"
)

type ISessionStore interface {
	Get(sessionCode string) (domain.Session, error)
	Create(frontSocketId string) (string, error)
	UpdateSocketId(sessionCode string, socketId string) error
	Delete(sessionCode string) error
	GetPlayers(sessionCode string) ([]domain.Player, error)
	AddPlayer(sessionCode string, player domain.Player) error
	RemovePlayer(sessionCode string, playerId uuid.UUID) error
}

type SessionStoreInMemory struct {
	generator tools.CodeGenerator
	sessions []domain.Session
}

//func (s *SessionStoreInMemory) GetPlayer(sessionCode string, playerId uuid.UUID) (domain.Player, error) {
//	session, err := s.Get(sessionCode)
//	if err != nil {
//		return domain.Player{}, err
//	}
//	for _, id := range session.PlayersIds {
//		if id == playerId {
//			return id, nil
//		}
//	}
//
//	return domain.Player{}, fmt.Errorf("SessionStoreInMemory: GetPlayer: player %s not found on session %s", playerId, sessionCode)
//}

func (s *SessionStoreInMemory) UpdateSocketId(sessionCode string, socketId string) error{
	session, err := s.Get(sessionCode)
	if err != nil {
		return err
	}
	session.FrontSocketId = socketId

	for index, sess := range s.sessions {
		if sess.SessionCode == sessionCode {
			s.sessions[index] = session
			return nil
		}
	}

	return fmt.Errorf("UpdateSocketId: Could not find session code: %s", sessionCode)
}

func (s *SessionStoreInMemory) GetPlayers(sessionCode string) ([]domain.Player, error) {
	session, err := s.Get(sessionCode)
	if err != nil {
		return nil, err
	}
	return session.Players, nil
}

func (s *SessionStoreInMemory) AddPlayer(sessionCode string, player domain.Player) error {

	for index, session := range s.sessions {
		if session.SessionCode == sessionCode {
			s.sessions[index].Players = append(s.sessions[index].Players, player)
			return nil
		}
	}

	return fmt.Errorf("session %s not found", sessionCode)
}

func removePlayer(session *domain.Session, playerId uuid.UUID) error {
	for index, player := range session.Players {
		if player["id"] == playerId {
			next := index + 1
			session.Players = append(session.Players[:index], session.Players[next:]...)
			return nil
		}
	}
	return fmt.Errorf("SessionStoreInMemory: RemovePlayer: player %s not found for session %s", playerId, session.SessionCode)
}

func (s *SessionStoreInMemory) RemovePlayer(sessionCode string, id uuid.UUID) error {
	for index, session := range s.sessions {
		if session.SessionCode == sessionCode {
			return removePlayer(&s.sessions[index], id)
		}
	}

	return fmt.Errorf("SessionStoreInMemory: RemovePlayer: session %s not found", sessionCode)
}

func NewSessionStoreInMemory(generator tools.CodeGenerator) *SessionStoreInMemory {
	return &SessionStoreInMemory{
		generator: generator,
		sessions:  []domain.Session{},
	}
}

func (s *SessionStoreInMemory) Get(sessionCode string) (domain.Session, error) {
	for _, session := range s.sessions {
		if session.SessionCode == sessionCode {
			return session, nil
		}
	}

	return domain.Session{}, fmt.Errorf("SessionStoreInMemory: Get: session %s not found", sessionCode)
}

func (s *SessionStoreInMemory) Create(frontSocketId string) (string, error) {
	code := s.generator.Generate()

	s.sessions = append(s.sessions, domain.Session{
		FrontSocketId: frontSocketId,
		SessionCode:   code,
		Players:       []domain.Player{},
	})

	return code, nil
}

func (s *SessionStoreInMemory) Delete(sessionCode string) error {
	for index, session := range s.sessions {
		if session.SessionCode == sessionCode {
			next := index + 1
			s.sessions = append(s.sessions[:index], s.sessions[next:]...)
			return nil
		}
	}

	return fmt.Errorf("SessionStoreInMemory: session %s not found", sessionCode)
}