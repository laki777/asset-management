package boltfunc

import (
	"encoding/json"
	"fmt"
	"log"
	"path"
	"runtime"
	"time"

	"github.com/boltdb/bolt"
)

var db *bolt.DB
var OpenDB bool

// https://github.com/boltdb/bolt#using-buckets

func Open(boltpath string) error {
	var err error
	_, filename, _, _ := runtime.Caller(0) // get full path of this file
	dbfile := path.Join(path.Dir(filename), boltpath)
	config := &bolt.Options{Timeout: 1 * time.Second}
	db, err = bolt.Open(dbfile, 0600, config)
	if err != nil {
		log.Fatal(err)
	}
	OpenDB = true
	return nil
}

func Close() {
	OpenDB = false
	db.Close()
}

type Person struct {
	Username  string
	FirstName string
	LastName  string
	Hardware  string
	Software  string
}

func (p *Person) Save(bucket string) error {
	if !OpenDB {
		return fmt.Errorf("db must be opened before saving")
	}
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		enc, err := p.encode()
		if err != nil {
			return fmt.Errorf("could not encode Username %s: %s", p.Username, err)
		}
		err = bucket.Put([]byte(p.Username), enc)
		return err
	})
	return err
}

func (p *Person) encode() ([]byte, error) {
	enc, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func decode(data []byte) (*Person, error) {
	var p *Person
	err := json.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func GetPerson(user, bucket string) (*Person, error) {
	if !OpenDB {
		return nil, fmt.Errorf("db must be opened before saving")
	}
	var p *Person
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		k := []byte(user)
		p, err = decode(b.Get(k))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Could not get Username %s", user)
		return nil, err
	}
	return p, nil
}

func List(bucket string) {
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(bucket))
		if c == nil {
			return nil
		}
		c1 := c.Cursor()
		for k, v := c1.First(); k != nil; k, v = c1.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}
		return nil
	})
}

// ADMIN PERSONS

type AdminPerson struct {
	Username  string
	Password  string
	FirstName string
	LastName  string
	Role      string
}

func (p *AdminPerson) AdminSave(bucket string) error {
	if !OpenDB {
		return fmt.Errorf("db must be opened before saving")
	}
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		enc, err := p.adminencode()
		if err != nil {
			return fmt.Errorf("could not encode Username %s: %s", p.Username, err)
		}
		err = bucket.Put([]byte(p.Username), enc)
		return err
	})
	return err
}

func (p *AdminPerson) adminencode() ([]byte, error) {
	enc, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func admindecode(data []byte) (*AdminPerson, error) {
	var p *AdminPerson
	err := json.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func AdminGetPerson(user, bucket string) (*AdminPerson, error) {
	if !OpenDB {
		return nil, fmt.Errorf("db must be opened before saving")
	}
	var p *AdminPerson
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		k := []byte(user)
		p, err = admindecode(b.Get(k))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		// 	fmt.Printf("Could not get Admin Username %s", user)
		return nil, err
	}
	return p, nil
}

func AdminGetFirstPerson() (string, error) {
	if !OpenDB {
		return "", fmt.Errorf("db must be opened before saving")
	}
	var role string
	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("adminpeople"))
		if c == nil {
			role = "admin"
			return nil
		}
		c1 := c.Cursor()
		k, _ := c1.First()
		if k == nil {
			role = "admin"
		} else {
			role = "user"
		}
		return nil
	})
	return role, err
}

// SESSIONS

type Session struct {
	Username string
	Time     time.Time
}

func (p *Session) SessionSave(cookieValue string) error {
	if !OpenDB {
		return fmt.Errorf("db must be opened before saving")
	}
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("sessions"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		enc, err := p.sessionencode()
		if err != nil {
			return fmt.Errorf("could not encode Username %s: %s", p.Username, err)
		}
		err = bucket.Put([]byte(cookieValue), enc)
		return err
	})
	return err
}

func (p *Session) sessionencode() ([]byte, error) {
	enc, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func sessiondecode(data []byte) (*Session, error) {
	var p *Session
	err := json.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func SessionGetPerson(CookieValue string) (*Session, error) {
	if !OpenDB {
		return nil, fmt.Errorf("db must be opened before saving")
	}
	var p *Session
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte("sessions"))
		if b == nil {
			return nil
		}
		k := []byte(CookieValue)
		if k == nil {
			return nil
		}
		p, err = sessiondecode(b.Get(k))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		// 	fmt.Printf("Could not get session Username %s", user)
		return nil, err
	}
	return p, nil
}

func SessionDelete(CookieValue string) error {
	if !OpenDB {
		return fmt.Errorf("db must be opened before saving")
	}
	err := db.Update(func(tx *bolt.Tx) error {
		var err error
		bucket := tx.Bucket([]byte("sessions"))
		if bucket == nil {
			return nil
		}

		err = bucket.Delete([]byte(CookieValue))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
