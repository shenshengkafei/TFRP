//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package storage

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"gopkg.in/mgo.v2"
)

// BaseDataProvider is the base struc of all data providers
type BaseDataProvider struct {
	Database string
	Password string
}

// Insert inserts a doc into collection
func (baseDataProvider *BaseDataProvider) Insert(collectionName string, id interface{}, doc interface{}) error {
	// Get session
	session, err := baseDataProvider.getDocDBSession()
	if err != nil {
		return err
	}

	defer session.Close()

	// SetSafe changes the session safety mode.
	// If the safe parameter is nil, the session is put in unsafe mode, and writes become fire-and-forget,
	// without error checking. The unsafe mode is faster since operations won't hold on waiting for a confirmation.
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode.
	session.SetSafe(&mgo.Safe{})

	// get collection
	collection := session.DB(baseDataProvider.Database).C(collectionName)

	_, err = collection.Upsert(id, doc)

	return err
}

// Find returns a doc from collection
func (baseDataProvider *BaseDataProvider) Find(collectionName string, qurey interface{}, result interface{}) error {
	// Get session
	session, err := baseDataProvider.getDocDBSession()
	if err != nil {
		return err
	}

	defer session.Close()

	// SetSafe changes the session safety mode.
	// If the safe parameter is nil, the session is put in unsafe mode, and writes become fire-and-forget,
	// without error checking. The unsafe mode is faster since operations won't hold on waiting for a confirmation.
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode.
	session.SetSafe(&mgo.Safe{})

	// get collection
	collection := session.DB(baseDataProvider.Database).C(collectionName)

	err = collection.Find(qurey).One(result)

	return err
}

// Remove deletes a doc from collection
func (baseDataProvider *BaseDataProvider) Remove(collectionName string, qurey interface{}) error {
	// Get session
	session, err := baseDataProvider.getDocDBSession()
	if err != nil {
		return err
	}

	defer session.Close()

	// SetSafe changes the session safety mode.
	// If the safe parameter is nil, the session is put in unsafe mode, and writes become fire-and-forget,
	// without error checking. The unsafe mode is faster since operations won't hold on waiting for a confirmation.
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode.
	session.SetSafe(&mgo.Safe{})

	// get collection
	collection := session.DB(baseDataProvider.Database).C(collectionName)

	err = collection.Remove(qurey)

	return err
}

func (baseDataProvider *BaseDataProvider) getDocDBSession() (*mgo.Session, error) {
	// DialInfo holds options for establishing a session with a MongoDB cluster.
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{fmt.Sprintf("%s.documents.azure.com:10255", baseDataProvider.Database)}, // Get HOST + PORT
		Timeout:  60 * time.Second,
		Database: baseDataProvider.Database, // It can be anything
		Username: baseDataProvider.Database, // Username
		Password: baseDataProvider.Password, // PASSWORD
		DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{})
		},
	}

	return mgo.DialWithInfo(dialInfo)
}
