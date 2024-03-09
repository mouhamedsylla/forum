package services

import (
	"fmt"
	"forum/Api/models"
	"forum/Api/storage"
	"net/http"
	"time"

	uuid "github.com/gofrs/uuid/v5"
)

type SessionManager struct {
	Sessions map[string]*models.Session
	Storage  *storage.Storage
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		Sessions: make(map[string]*models.Session, 0),
		Storage:  storage.NewStorage(),
	}
}

func GenerateUUID() string {
	id, _ := uuid.NewV7()
	return id.String()
}

func (sm *SessionManager) CreateSession(user *models.User) *models.Session {
	s := &models.Session{
		SessionId: GenerateUUID(),
		User:      user,
		TimeOut:   time.Now().Add(4 * time.Minute),
	}
	sm.Sessions[s.SessionId] = s
	return s
}

func (sm *SessionManager) DeleteSessionWithUserID(userId int) {
	for _, session := range sm.Sessions {
		if session.User.Id == userId {
			sm.DeleteSession(session.SessionId)
		}
	}
}

func (sm *SessionManager) DeleteSession(token string) {
	delete(sm.Sessions, token)
}

func (sm *SessionManager) GetSession(idSession string) *models.SessionDb {
	result := sm.Storage.Select(models.SessionDb{}, "SessionID", idSession).([]models.SessionDb)
	if len(result) == 0 {
		return nil
	}
	return &result[0]
}

func (sm *SessionManager) GenerateSession(w *http.ResponseWriter, user *models.User) {
	session := sm.CreateSession(user)
	http.SetCookie(*w, &http.Cookie{
		Name:    "forum",
		Value:   session.SessionId,
		Expires: session.TimeOut,
	})
	sm.Storage.Insert(models.SessionDb{
		SessionID: session.SessionId,
		User_Id:   user.Id,
		TimeOut:   session.TimeOut.Format("2006-01-02 15:04:05"),
	})
}

func IsExpired(Time string) bool {
	TimeOut, err := time.Parse("2006-01-02 15:04:05", Time)
	if err != nil {
		fmt.Println(err)
	}
	return TimeOut.Before(time.Now())
}

func (sm *SessionManager) InitSessionMap() {
	sessionsDb := sm.Storage.SelectAll(models.SessionDb{}).([]models.SessionDb)
	if len(sessionsDb) != 0 {
		for _, sDb := range sessionsDb {
			if IsExpired(sDb.TimeOut) {
				sm.Storage.DeleteExipireSession(sDb.User_Id)
				continue
			}
			user := sm.Storage.Select(models.User{}, "Id", sDb.User_Id).([]models.User)[0]
			t, _ := time.Parse("2006-01-02 15:04:05", sDb.TimeOut)
			sm.Sessions[sDb.SessionID] = &models.Session{
				SessionId: sDb.SessionID,
				User: &user,
				TimeOut: t,
			}
		}
	}
}

var Global_SessionManager = NewSessionManager()
