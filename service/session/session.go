package session

import (
    "net/http"

    "github.com/gorilla/sessions"
)

type Manager struct {
    store *sessions.CookieStore
}

func NewManager(secretKey []byte) *Manager {
    sessionStore := sessions.NewCookieStore(secretKey)

    return &Manager{store: sessionStore}
}

// حفظ قيمة داخل جلسة
func (m *Manager) SetValue(w http.ResponseWriter, r *http.Request, sessionName string, key string, value interface{}) error {
    session, _ := m.store.Get(r, sessionName)
    session.Values[key] = value
    return session.Save(r, w)
}

// الحصول على قيمة من الجلسة
func (m *Manager) GetValue(r *http.Request, sessionName string, key string) (interface{}, bool) {
    session, _ := m.store.Get(r, sessionName)
    val, ok := session.Values[key]
    return val, ok
}

// حذف الجلسة
func (m *Manager) ClearSession(w http.ResponseWriter, r *http.Request, sessionName string) error {
    session, _ := m.store.Get(r, sessionName)
    session.Options.MaxAge = -1
    return session.Save(r, w)
}
