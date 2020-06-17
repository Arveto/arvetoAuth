// Copyright (c) 2020, HuguesGuilleus. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/prologic/bitcask"
	"log"
	"reflect"
)

var (
	// The default option
	opt = []bitcask.Option{
		bitcask.WithSync(false),
		bitcask.WithMaxValueSize(1000000),
	}
	endFind = errors.New("end find")
)

// DB is a data base
type DB struct {
	intern   *bitcask.Bitcask
	maxIndex Key // The max id
}

// Open a new DabatBase
func New(name string) (db *DB) {
	intern, err := bitcask.Open(name, opt...)
	if err != nil {
		log.Println(err)
	}

	db = &DB{
		intern: intern,
	}

	db.intern.Fold(func(key []byte) error {
		if k := keyBytes(key); k > db.maxIndex {
			db.maxIndex = k
		}
		return nil
	})

	return
}

/* BASIC MANIPULATION */

// Create a new index in the DB.
func (db *DB) New() Key {
	db.maxIndex++
	return db.maxIndex - 1
}

// Return true if the key is unknow in the DB.
func (db *DB) Unknown(key Key) bool {
	return !db.intern.Has(key.bytes())
}
func (db *DB) UnknownS(key string) bool {
	return !db.intern.Has([]byte(key))
}

// Delete on element
func (db *DB) Delete(key Key) {
	if err := db.intern.Delete(key.bytes()); err != nil {
		log.Println(err)
	}
}

// Delete on element
func (db *DB) DeleteS(key string) {
	if err := db.intern.Delete([]byte(key)); err != nil {
		log.Println(err)
	}
}

// Remove all keys in the DB.
func (db *DB) DeleteAll() {
	db.maxIndex = 0
	if err := db.intern.DeleteAll(); err != nil {
		log.Println(err)
	}
}

/* GET */

// Get the elemena andd sav it in v.
func (db *DB) Get(key Key, v interface{}) (noExist bool) {
	data, err := db.intern.Get(key.bytes())
	return get(data, err, v)
}

func (db *DB) GetRaw(key Key) []byte {
	data, err := db.intern.Get(key.bytes())

	if err != nil && err != bitcask.ErrKeyNotFound {
		log.Println(err)
		return nil
	}

	return data
}
func (db *DB) GetSRaw(key string) []byte {
	data, err := db.intern.Get([]byte(key))

	if err != nil && err != bitcask.ErrKeyNotFound {
		log.Println(err)
		return nil
	}

	return data
}

// Get the element with a string key.
func (db *DB) GetS(key string, v interface{}) (noExist bool) {
	data, err := db.intern.Get([]byte(key))
	return get(data, err, v)
}

func get(data []byte, err error, v interface{}) (noExist bool) {
	// Set v to Zero
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))

	// Checvk the error
	if err == bitcask.ErrKeyNotFound {
		return true
	} else if err != nil {
		log.Println(err)
		return true
	}

	// Decode the value
	err = gob.NewDecoder(bytes.NewReader(data)).Decode(v)
	if err != nil {
		log.Println(err)
	}

	return false
}

// Make an iteration on all the element in the DB. It must be a function
// that take a Key and a other type for the value, else it panic.
//	MyDB.For(func(k db.Key, v MyType){...})
func (db *DB) For(it interface{}) {
	f := reflect.ValueOf(it)
	t := f.Type()
	if t.Kind() != reflect.Func ||
		t.NumIn() != 2 ||
		t.In(0) != reflect.TypeOf(Key(0)) {
		log.Panic("DB.For() need a iteration function")
	}
	v := reflect.New(t.In(1)).Elem()

	db.intern.Fold(func(key []byte) error {
		data, err := db.intern.Get(key)
		if err != nil {
			log.Println(err)
			return nil
		}

		v.Set(reflect.Zero(v.Type()))
		err = gob.NewDecoder(bytes.NewReader(data)).DecodeValue(v)
		if err != nil {
			log.Println(err)
			return nil
		}

		f.Call([]reflect.Value{reflect.ValueOf(keyBytes(key)), v})

		return nil
	})
}

// Make an iteration on all the element in the DB that begins with prefix.
// It must be a function that take a string and a other type for the value,
// else it panic.
//	MyDB.For(func(k db.Key, v MyType){...})
func (db *DB) ForS(prefix string, it interface{}) {
	f := reflect.ValueOf(it)
	t := f.Type()
	if t.Kind() != reflect.Func ||
		t.NumIn() != 2 ||
		t.In(0) != reflect.TypeOf("") {
		log.Panic("DB.ForS() need a iteration function")
	}
	v := reflect.New(t.In(1)).Elem()

	db.intern.Scan([]byte(prefix), func(key []byte) error {
		data, err := db.intern.Get(key)
		if err != nil {
			log.Println(err)
			return nil
		}

		v.Set(reflect.Zero(v.Type()))
		err = gob.NewDecoder(bytes.NewReader(data)).DecodeValue(v)
		if err != nil {
			log.Println(err)
			return nil
		}

		f.Call([]reflect.Value{reflect.ValueOf(string(key)), v})

		return nil
	})
}

/* SET */

func (db *DB) Set(key Key, v interface{}) {
	db.set(key.bytes(), v)
}

func (db *DB) SetS(key string, v interface{}) {
	db.set([]byte(key), v)
}

func (db *DB) set(k []byte, v interface{}) {
	w := bytes.Buffer{}
	if err := gob.NewEncoder(&w).Encode(v); err != nil {
		log.Println(err)
		return
	}

	if err := db.intern.Put(k, w.Bytes()); err != nil {
		log.Println(err)
	}
}

func (db *DB) SetRaw(key Key, raw []byte) {
	if err := db.intern.Put(key.bytes(), raw); err != nil {
		log.Println(err)
	}
}

func (db *DB) SetSRaw(key string, raw []byte) {
	if err := db.intern.Put([]byte(key), raw); err != nil {
		log.Println(err)
	}
}
