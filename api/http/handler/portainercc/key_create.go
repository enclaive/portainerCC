package portainercc

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"os/exec"
	"sync"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/dataservices"
)

type KeyParams struct {
	KeyType            string
	Description        string
	TeamAccessPolicies portainer.TeamAccessPolicies
	Data               string
}

type KeyResponse struct {
	Id                 portainer.KeyID
	KeyType            string
	Description        string
	TeamAccessPolicies portainer.TeamAccessPolicies
}

func (handler *Handler) createKey(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var params KeyParams
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		return httperror.BadRequest("request body malefomred", err)
	}

	keyObject, errx := CreateNewKey(handler.DataStore.Key(), params.KeyType, params.Description, params.TeamAccessPolicies, params.Data)

	if errx != nil {
		return errx
	}

	res := KeyResponse{
		Id:                 keyObject.ID,
		KeyType:            keyObject.KeyType,
		Description:        keyObject.Description,
		TeamAccessPolicies: keyObject.TeamAccessPolicies,
	}

	return response.JSON(w, res)
}

func CreateNewKey(keyService dataservices.KeyService, keyType string, description string, access portainer.TeamAccessPolicies, importData string) (*portainer.Key, *httperror.HandlerError) {
	keyObject := &portainer.Key{
		KeyType:            keyType,
		Description:        description,
		TeamAccessPolicies: access,
	}

	//handle key type
	if keyObject.KeyType == "SIGNING" {
		if importData != "" {
			//IMPORT data as key
			block, _ := pem.Decode([]byte(importData))
			if block == nil {
				return nil, httperror.InternalServerError("could not import rsa key, invalid pem", nil)
			}

			privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if block == nil {
				return nil, httperror.InternalServerError("could not import rsa key", err)
			}

			keyObject.SigningKey = privKey
		} else {
			//gen new key
			privKey, err := GenerateMultiPrimeKeyE3(rand.Reader, 2, 3072)
			if err != nil {
				return nil, httperror.InternalServerError("could not generate rsa key", err)
			}

			keyObject.SigningKey = privKey
		}
	} else if keyObject.KeyType == "FILE_ENC" {
		//gramine pf file key
		if importData != "" {
			//IMPORT data as key - TODO maybe hex check/size
			keyObject.PFKey = importData
		} else {
			//gen new key
			tempKeyFile, err := ioutil.TempFile("", "super")
			if err != nil {
				return nil, httperror.InternalServerError("could not generate file key", err)
			}
			// defer os.Remove(tempKeyFile.Name())

			//create key with sgx
			cmd := exec.Command("gramine-sgx-pf-crypt", "gen-key", "-w", tempKeyFile.Name())
			_, err = cmd.Output()
			if err != nil {
				return nil, httperror.InternalServerError("could not generate file key", err)
			}

			//save as base64
			file, err := ioutil.ReadFile(tempKeyFile.Name())
			if err != nil {
				return nil, httperror.InternalServerError("could not generate file key", err)
			}

			keyObject.PFKey = base64.StdEncoding.EncodeToString(file)
		}
	} else {
		return nil, httperror.InternalServerError("invalid key type", nil)
	}

	err := keyService.Create(keyObject)

	if err != nil {
		return nil, httperror.InternalServerError("could save key in db", err)
	}

	return keyObject, nil
}

///// fixed public exponent of 3 for sgx signing key

var bigOne = big.NewInt(1)

// GenerateMultiPrimeKey generates a multi-prime RSA keypair of the given bit
// size and the given random source, as suggested in [1]. Although the public
// keys are compatible (actually, indistinguishable) from the 2-prime case,
// the private keys are not. Thus it may not be possible to export multi-prime
// private keys in certain formats or to subsequently import them into other
// code.
//
// Table 1 in [2] suggests maximum numbers of primes for a given size.
//
// [1] US patent 4405829 (1972, expired)
// [2] http://www.cacr.math.uwaterloo.ca/techreports/2006/cacr2006-16.pdf
func GenerateMultiPrimeKeyE3(random io.Reader, nprimes int, bits int) (*rsa.PrivateKey, error) {
	MaybeReadByte(random)

	priv := new(rsa.PrivateKey)
	priv.E = 3

	if nprimes < 2 {
		return nil, errors.New("crypto/rsa: GenerateMultiPrimeKey: nprimes must be >= 2")
	}

	if bits < 64 {
		primeLimit := float64(uint64(1) << uint(bits/nprimes))
		// pi approximates the number of primes less than primeLimit
		pi := primeLimit / (math.Log(primeLimit) - 1)
		// Generated primes start with 11 (in binary) so we can only
		// use a quarter of them.
		pi /= 4
		// Use a factor of two to ensure that key generation terminates
		// in a reasonable amount of time.
		pi /= 2
		if pi <= float64(nprimes) {
			return nil, errors.New("crypto/rsa: too few primes of given length to generate an RSA key")
		}
	}

	primes := make([]*big.Int, nprimes)

NextSetOfPrimes:
	for {
		todo := bits
		// crypto/rand should set the top two bits in each prime.
		// Thus each prime has the form
		//   p_i = 2^bitlen(p_i) × 0.11... (in base 2).
		// And the product is:
		//   P = 2^todo × α
		// where α is the product of nprimes numbers of the form 0.11...
		//
		// If α < 1/2 (which can happen for nprimes > 2), we need to
		// shift todo to compensate for lost bits: the mean value of 0.11...
		// is 7/8, so todo + shift - nprimes * log2(7/8) ~= bits - 1/2
		// will give good results.
		if nprimes >= 7 {
			todo += (nprimes - 2) / 5
		}
		for i := 0; i < nprimes; i++ {
			var err error
			primes[i], err = rand.Prime(random, todo/(nprimes-i))
			if err != nil {
				return nil, err
			}
			todo -= primes[i].BitLen()
		}

		// Make sure that primes is pairwise unequal.
		for i, prime := range primes {
			for j := 0; j < i; j++ {
				if prime.Cmp(primes[j]) == 0 {
					continue NextSetOfPrimes
				}
			}
		}

		n := new(big.Int).Set(bigOne)
		totient := new(big.Int).Set(bigOne)
		pminus1 := new(big.Int)
		for _, prime := range primes {
			n.Mul(n, prime)
			pminus1.Sub(prime, bigOne)
			totient.Mul(totient, pminus1)
		}
		if n.BitLen() != bits {
			// This should never happen for nprimes == 2 because
			// crypto/rand should set the top two bits in each prime.
			// For nprimes > 2 we hope it does not happen often.
			continue NextSetOfPrimes
		}

		priv.D = new(big.Int)
		e := big.NewInt(int64(priv.E))
		ok := priv.D.ModInverse(e, totient)

		if ok != nil {
			priv.Primes = primes
			priv.N = n
			break
		}
	}

	priv.Precompute()
	return priv, nil
}

// MaybeReadByte reads a single byte from r with ~50% probability. This is used
// to ensure that callers do not depend on non-guaranteed behaviour, e.g.
// assuming that rsa.GenerateKey is deterministic w.r.t. a given random stream.
//
// This does not affect tests that pass a stream of fixed bytes as the random
// source (e.g. a zeroReader).

var (
	closedChanOnce sync.Once
	closedChan     chan struct{}
)

func MaybeReadByte(r io.Reader) {
	closedChanOnce.Do(func() {
		closedChan = make(chan struct{})
		close(closedChan)
	})

	select {
	case <-closedChan:
		return
	case <-closedChan:
		var buf [1]byte
		r.Read(buf[:])
	}
}
