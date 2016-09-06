//
// 3nigm4 auth package
// Author: Guido Ronchetti <dyst0ni3@gmail.com>
// v1.0 16/06/2016
//
// This mock database is used for tests purposes, should
// never be used in production environment. It's not
// concurrency safe and do not implement any performance
// optimisation logic.
//

package dbmock

// Golang std libs
import (
	"encoding/hex"
	"fmt"
)

import (
	ty "github.com/nexocrew/3nigm4/lib/auth/types"
	db "github.com/nexocrew/3nigm4/lib/database/client"
)

type Mockdb struct {
	addresses string
	user      string
	password  string
	authDb    string
	// in memory storage
	userStorage    map[string]*ty.User
	sessionStorage map[string]*ty.Session
}

func NewMockDb(args *db.DbArgs) *Mockdb {
	return &Mockdb{
		addresses:      composeDbAddress(args),
		user:           args.User,
		password:       args.Password,
		authDb:         args.AuthDb,
		userStorage:    make(map[string]*ty.User),
		sessionStorage: make(map[string]*ty.Session),
	}
}

func (d *Mockdb) Copy() db.Database {
	return d
}

func (d *Mockdb) Close() {
}

func (d *Mockdb) GetUser(username string) (*ty.User, error) {
	user, ok := d.userStorage[username]
	if !ok {
		return nil, fmt.Errorf("unable to find the required %s user", username)
	}
	return user, nil
}

func (d *Mockdb) SetUser(user *ty.User) error {
	_, ok := d.userStorage[user.Username]
	if ok {
		return fmt.Errorf("user %s already exist in the db", user.Username)
	}
	d.userStorage[user.Username] = user
	return nil
}

func (d *Mockdb) RemoveUser(username string) error {
	if _, ok := d.userStorage[username]; !ok {
		return fmt.Errorf("unable to find required %s user", username)
	}
	delete(d.userStorage, username)
	return nil
}

func (d *Mockdb) GetSession(token []byte) (*ty.Session, error) {
	h := hex.EncodeToString(token)
	session, ok := d.sessionStorage[h]
	if !ok {
		return nil, fmt.Errorf("unable to find the required %s session", h)
	}
	return session, nil
}

func (d *Mockdb) SetSession(s *ty.Session) error {
	h := hex.EncodeToString(s.Token)
	d.sessionStorage[h] = s
	return nil
}

func (d *Mockdb) RemoveSession(token []byte) error {
	h := hex.EncodeToString(token)
	if _, ok := d.sessionStorage[h]; !ok {
		return fmt.Errorf("unable to find required %s session", h)
	}
	delete(d.sessionStorage, h)
	return nil
}

func (d *Mockdb) RemoveAllSessions() error {
	d.sessionStorage = make(map[string]*ty.Session)
	return nil
}

// composeDbAddress compose a string starting from dbArgs slice.
func composeDbAddress(args *db.DbArgs) string {
	dbAccess := fmt.Sprintf("mongodb://%s:%s@", args.User, args.Password)
	for idx, addr := range args.Addresses {
		dbAccess += addr
		if idx != len(args.Addresses)-1 {
			dbAccess += ","
		}
	}
	dbAccess += fmt.Sprintf("/?authSource=%s", args.AuthDb)
	return dbAccess
}
