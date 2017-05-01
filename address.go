package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/boltdb/bolt"

	yaml "gopkg.in/yaml.v2"
)

type Addresses map[string]*Address
type AddrByTown map[string][]*Address

type Address struct {
	Name    string
	Town    string
	Address string
	Comment string
	Type    string
}

var addresses Addresses

func Add(a *Address) error {
	if addresses == nil {
		addresses = make(map[string]*Address)
	}

	if _, ok := addresses[a.Key()]; ok {
		return errors.New("This Town/Name combinasion exists already.")
	}

	addresses[a.Key()] = a
	return nil
}

func (a *Address) Key() string {
	return a.Town + "/" + a.Name
}

func (a *Address) Yaml() string {
	out, _ := yaml.Marshal(*a)
	return string(out)
}

func ParseAddresses(input []byte) (as []*Address, err error) {
	err = yaml.Unmarshal(input, &as)

	return
}

func AddFromDB() {
}

func SaveToDB() {
	err := db.Update(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte("address"))
		if b == nil {
			b, err = tx.CreateBucket([]byte("address"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
		}

		for _, a := range addresses {
			if v := b.Get([]byte(a.Key())); v != nil {
				// Key exists
				fmt.Println(a.Key(), "was not saved in DB because it already exists.")
			}
			err = b.Put([]byte(a.Key()), []byte(a.Yaml()))
			if err != nil {
				fmt.Println("Error saving", a.Key(), "to DB.")
			}
		}

		return nil
	})
	if err != nil {
		fmt.Println("Transaction failed:", err)
	}
}

func AddFromYaml(input string) error {
	af, err := os.Open(input)
	if err != nil {
		return err
	}
	ab, err := ioutil.ReadAll(af)
	if err != nil {
		return err
	}

	as, err := ParseAddresses(ab)
	if err != nil {
		return err
	}

	for _, a := range as {
		err := Add(a)
		if err != nil {
			fmt.Println(a, "could not be added:", err)
		}
	}

	return nil
}

func ByTown() AddrByTown {
	ret := make(AddrByTown)

	for _, a := range addresses {
		ret[a.Town] = append(ret[a.Town], a)
	}

	return ret
}
