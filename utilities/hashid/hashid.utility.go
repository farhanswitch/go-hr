package hashidutility

import (
	"log"
	"os"

	hashids "github.com/speps/go-hashids/v2"
)

var hash *hashids.HashID

func FactoryHashID() *hashids.HashID {
	if hash == nil {
		hd := hashids.NewData()
		salt := os.Getenv("HASHID_SALT")
		if salt == "" {
			log.Fatalf("HashID Salt is undefined!")

		}
		hd.Salt = salt
		hd.MinLength = 6
		hash, _ = hashids.NewWithData(hd)
	}
	return hash
}
