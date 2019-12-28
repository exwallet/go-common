/*
 * @Author: kidd
 * @Date: 12/24/19 3:53 PM
 */

package sessionManager

import (
	"fmt"
	"github.com/exwallet/go-common/cache/redis"
	"github.com/exwallet/go-common/goutil/gostring"
	"github.com/exwallet/go-common/goutil/gotime"
	"github.com/exwallet/go-common/log"
	"sync"
)

type cookieKey struct {
	SessionId  string
	UserId     string
	UserStatus string
	Lan        string
	Role       string
}

type SessionManager struct {
	env             string       // 前缀
	CookieKey       *cookieKey   //
	saveLifeSeconds int64        // session 执行续命时间间隔
	survivalSeconds int64        // session 存活时间
	maxSaveLifeTime int          // 最大续命次数
	lock            sync.RWMutex //
}

// save to cache
type Session struct {
	Sid              string //
	UserId           int64  //
	Username         string //
	LastActivityTime int64  //
	LoginIp          string //
	UserStatus       int64  //
	Lan              string //
	GoogleSecret     string //
	RoleId           int64  //
	SaveLifeTime     int    //
}

func NewSessManager(env string, saveLifeSeconds int64, maxSaveLifeTime int, survivalSeconds int64) *SessionManager {
	return &SessionManager{
		env: env,
		CookieKey: &cookieKey{
			SessionId:  env + "SessionID",
			UserId:     env + "UserId",
			UserStatus: env + "Code",
			Lan:        env + "Lan",
			Role:       env + "Role",
		},
		saveLifeSeconds: saveLifeSeconds,
		survivalSeconds: survivalSeconds,
		maxSaveLifeTime: maxSaveLifeTime,
	}
}

func (m *SessionManager) cacheKeySid(sid string) string {
	return fmt.Sprintf("%s_sid_%s", m.env, sid)
}

func (m *SessionManager) cacheKeyUid2Sid(uid int64) string {
	return fmt.Sprintf("%s_uid2sid_%d", m.env, uid)
}
func (m *SessionManager) cacheKeySaveLifeTime(uid int64) string {
	return fmt.Sprintf("%s_saveLifeTime_%d", m.env, uid)
}

func (m *SessionManager) _setCookie(ctx ContextInf, key string, val string) {
	//if key == m.CookieKey.SessionId {
	//	ctx.SetCookie(key, val, m.survivalSeconds, "/", "", false, true)
	//} else {
	//	ctx.SetCookie(key, val, m.survivalSeconds, "/", "", false, false)
	//}
	ctx.SetCookie(key, val, m.survivalSeconds, "/", "", false, false)
}

func (m *SessionManager) _expiredCookie(ctx ContextInf, key string) {
	ctx.SetCookie(key, "", -1, "/", ctx.GetHeader("Domain"), false, true)
}

func (m *SessionManager) _updateCookies(ctx ContextInf, s *Session) {
	m._setCookie(ctx, m.CookieKey.SessionId, s.Sid)
	m._setCookie(ctx, m.CookieKey.UserId, fmt.Sprintf("%d", s.UserId))
	m._setCookie(ctx, m.CookieKey.UserStatus, fmt.Sprintf("%d", s.UserStatus))
	m._setCookie(ctx, m.CookieKey.Lan, s.Lan)
	m._setCookie(ctx, m.CookieKey.Role, fmt.Sprintf("%d", s.RoleId))
}

func (m *SessionManager) InitSessionId(ctx ContextInf, keepIfNull ...bool) string {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m._initSessionId(ctx, keepIfNull...)
}
func (m *SessionManager) _initSessionId(ctx ContextInf, keepIfNull ...bool) string {
	var sid string
	sid = ctx.GetCookie(m.CookieKey.SessionId)
	if len(sid) < 32 {
		sid = ""
	}
	if sid == "" {
		if len(keepIfNull) > 0 && keepIfNull[0] {
			return sid
		}
		sid = gostring.GenerateGuid()
		log.Info("IP[%s]取得新SessionId[%s]", ctx.IP(), sid)
		m._setCookie(ctx, m.CookieKey.SessionId, sid)
		return sid
	}
	log.Debug("IP[%s]使用旧SessionId[%s]", ctx.IP(), sid)
	return sid
}

func (m *SessionManager) GetSession(ctx ContextInf) *Session {
	m.lock.Lock()
	defer m.lock.Unlock()
	sid := ctx.GetCookie(m.CookieKey.SessionId)
	if sid == "" {
		return nil
	}
	s := m._getSessionBySid(sid)
	if s != nil {
		// 续命
		if m.maxSaveLifeTime > 0 {
			if gotime.UnixNowMillSec()-s.LastActivityTime > m.saveLifeSeconds*1000 {
				if s.SaveLifeTime < m.maxSaveLifeTime {
					log.Info("Session: Evn[%s]user[%v][%s] 最后活跃时间[%s] 执行session续命",
						m.env, s.UserId, s.Username, gotime.MillSecStrCST8(s.LastActivityTime))
					s.LastActivityTime = gotime.UnixNowMillSec()
					s.SaveLifeTime += 1
					m._updateCookies(ctx, s)
					m._updateSession(s)
				}
			}
		}
		return s
	}
	return nil
}

func (m *SessionManager) GetSessionBySid(sid string) *Session {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m._getSessionBySid(sid)
}
func (m *SessionManager) _getSessionBySid(sid string) *Session {
	k := m.cacheKeySid(sid)
	obj, _ := redis.GetObj(k, (*Session)(nil))
	if s, ok := obj.(*Session); ok {
		return s
	}
	return nil
}

func (m *SessionManager) GetSessionByUid(uid int64) *Session {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m._getSessionByUid(uid)
}
func (m *SessionManager) _getSessionByUid(uid int64) *Session {
	sid, _ := redis.Get(m.cacheKeyUid2Sid(uid))
	if sid == "" {
		return nil
	}
	return m._getSessionBySid(sid)
}

func (m *SessionManager) DoLogin(ctx ContextInf, userId int64, username string, userStatus int64, roleId int64, googleSecret string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	s := &Session{
		Sid:              m._initSessionId(ctx),
		UserId:           userId,
		Username:         username,
		LastActivityTime: gotime.UnixNowMillSec(),
		LoginIp:          ctx.IP(),
		UserStatus:       userStatus,
		Lan:              ctx.GetCookie(m.CookieKey.Lan),
		GoogleSecret:     googleSecret,
		RoleId:           roleId,
	}
	if s.Lan == "" {
		s.Lan = "en"
	}
	m._updateCookies(ctx, s)
	m._updateSession(s)
}

func (m *SessionManager) DoLogout(ctx ContextInf) {
	m.lock.Lock()
	defer m.lock.Unlock()
	log.Info("清退Session: Env[%s]用户[%s]", m.env, m.CookieKey.UserId)
	sid := ctx.GetCookie(m.CookieKey.SessionId)
	if sid != "" {
		s := m._getSessionBySid(sid)
		if s != nil {
			m._removeSession(s)
		}
	}
	m._expiredCookie(ctx, m.CookieKey.SessionId)
	m._expiredCookie(ctx, m.CookieKey.UserId)
	m._expiredCookie(ctx, m.CookieKey.UserStatus)
	m._expiredCookie(ctx, m.CookieKey.Lan)
	m._expiredCookie(ctx, m.CookieKey.Role)
}

func (m *SessionManager) DoLogoutByUserId(uid int64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	log.Info("清退Session: Env[%s]用户[%v]", m.env, uid)
	s := m._getSessionByUid(uid)
	if s != nil {
		m._removeSession(s)
	}
}

func (m *SessionManager) UpdateSession(s *Session) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m._updateSession(s)
}
func (m *SessionManager) _updateSession(s *Session) {
	redis.SetObjAndExpire(m.cacheKeySid(s.Sid), s, int(m.survivalSeconds))
	redis.SetAndExpire(m.cacheKeyUid2Sid(s.UserId), s.Sid, int(m.survivalSeconds))
	log.Info("保存Session: Env[%s]用户[%v]session", m.env, s.UserId)
}

func (m *SessionManager) _removeSession(s *Session) {
	redis.Delete(m.cacheKeyUid2Sid(s.UserId)) // 去掉 userId - sessionId 映射关系
	redis.Delete(m.cacheKeySid(s.Sid))        // delete session
	log.Warn("清除Session: Env[%s]用户[%v]session", m.env, s.UserId)
}
