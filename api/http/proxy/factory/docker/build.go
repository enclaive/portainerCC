package docker

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"strconv"
	"sync"

	portainer "github.com/portainer/portainer/api"

	"github.com/portainer/portainer/api/archive"

	"github.com/rs/zerolog/log"
)

const OneMegabyte = 1024768

type postDockerfileRequest struct {
	Content string
}

// buildOperation inspects the "Content-Type" header to determine if it needs to alter the request.
// If the value of the header is empty, it means that a Dockerfile is posted via upload, the function
// will extract the file content from the request body, tar it, and rewrite the body.
// If the value of the header contains "application/json", it means that the content of a Dockerfile is posted
// in the request payload as JSON, the function will create a new file called Dockerfile inside a tar archive and
// rewrite the body of the request.
// In any other case, it will leave the request unaltered.
func buildOperation(request *http.Request, transport *Transport) error {
	contentTypeHeader := request.Header.Get("Content-Type")

	mediaType, _, err := mime.ParseMediaType(contentTypeHeader)
	if err != nil {
		return err
	}

	var buffer []byte
	switch mediaType {
	case "":
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return err
		}

		buffer, err = archive.TarFileInBuffer(body, "Dockerfile", 0600)
		if err != nil {
			return err
		}

	case "application/json":
		var req postDockerfileRequest
		if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
			return err
		}

		buffer, err = archive.TarFileInBuffer([]byte(req.Content), "Dockerfile", 0600)
		if err != nil {
			return err
		}

	case "multipart/form-data":
		err := request.ParseMultipartForm(32 * OneMegabyte) // limit parser memory to 32MB
		if err != nil {
			return err
		}

		if request.MultipartForm == nil || request.MultipartForm.File == nil {
			return errors.New("uploaded files not found to build image")
		}

		tfb := archive.NewTarFileInBuffer()
		defer tfb.Close()

		for k := range request.MultipartForm.File {
			f, hdr, err := request.FormFile(k)
			if err != nil {
				return err
			}

			defer f.Close()

			log.Info().Str("filename", hdr.Filename).Int64("size", hdr.Size).Msg("upload the file to build image")

			content, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}

			filename := hdr.Filename
			if hdr.Filename == "blob" {
				filename = "Dockerfile"
			}

			if err := tfb.Put(content, filename, 0600); err != nil {
				return err
			}
		}

		buffer = tfb.Bytes()
		request.Form = nil
		request.PostForm = nil
		request.MultipartForm = nil

	default:
		return nil
	}

	request.Body = ioutil.NopCloser(bytes.NewReader(buffer))
	request.ContentLength = int64(len(buffer))
	request.Header.Set("Content-Type", "application/x-tar")

	//check for sgxbuild param
	//check if buildWithSGX is set
	sgx := request.URL.Query().Get("sgx")
	if sgx != "true" {
		return nil
	}

	signingKeyIdStr := request.URL.Query().Get("signing-key-id")
	if signingKeyIdStr == "" {
		return errors.New("missing signing key id")
	}

	signingKeyId, err := strconv.ParseInt(signingKeyIdStr, 10, 64)
	if err != nil {
		return errors.New("singingkeyid is not an int")
	}

	log.Info().Str("SGX Value", sgx).Int64("SigninKeyID", signingKeyId).Msg("Build With SGX docker subcontainer")

	//read the signing key out of db
	var keyObject *portainer.Key
	keyObject, err = transport.dataStore.Key().Key(portainer.KeyID(signingKeyId))
	if err != nil {
		return errors.New("failed to retreive the signing key from the data store")
	}

	// to pem
	privKeyBytes := x509.MarshalPKCS1PrivateKey(keyObject.SigningKey)
	signingKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privKeyBytes,
		},
	)

	//handle subrequests for subcontainer
	type hostConfig struct {
		Binds []string
	}

	type createContainerValues struct {
		Image      string
		Cmd        []string
		HostConfig hostConfig
	}

	type createContainerResponse struct {
		Id       string
		Warnings []string
	}

	//create docker container
	// TODO docker image muss da sein
	hostCfg := hostConfig{
		Binds: []string{
			"/var/run/docker.sock:/var/run/docker.sock:z",
			"/var/run/docker.sock:/var/run/alternative.sock:z",
			"/tmp/docker-context:/tmp/docker-context",
			"testVolume:/testVOLUUUUUME",
		},
	}

	postBody, _ := json.Marshal(createContainerValues{
		"docker",
		[]string{"sh", "-c", fmt.Sprintf("echo \"%s\" > /tmp/secret && cd && tar xf /tmp/docker-context && BUILDKIT_PROGRESS=plain DOCKER_BUILDKIT=1 docker build --no-cache --secret id=PORTAINER_SGX_SIGNER_KEY,src=/tmp/secret -t freaking .", signingKeyPEM)},
		hostCfg,
	})

	subRequest := request.Clone(request.Context())

	dockerContext, err := io.ReadAll(subRequest.Body)
	_ = os.WriteFile("/tmp/docker-context", dockerContext, 0644)

	subRequest.URL.Path = "/containers/create"
	subRequest.Body = ioutil.NopCloser(bytes.NewReader(postBody))
	subRequest.ContentLength = int64(len(postBody))
	subRequest.Header.Set("Content-Type", "application/json")

	resp, err := transport.executeDockerRequest(subRequest)

	body, err := ioutil.ReadAll(resp.Body)
	var createResponse createContainerResponse
	json.Unmarshal(body, &createResponse)

	log.Info().Str("Container ID", createResponse.Id).Msg("Created docker subcontainer")

	// run docker container
	subRequest.URL.Path = fmt.Sprintf("/containers/%s/start", createResponse.Id)
	subRequest.Body = nil
	subRequest.ContentLength = 0

	resp, err = transport.executeDockerRequest(subRequest)

	log.Info().Str("Container ID", createResponse.Id).Msg("Started docker subcontainer, waiting to finish")

	// wait for container do finish
	subRequest.URL.Path = fmt.Sprintf("/containers/%s/wait", createResponse.Id)

	//we need to wait here
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		resp, err = transport.executeDockerRequest(subRequest)
	}()

	wg.Wait()

	fmt.Println(resp)

	//response to user will be stdout form subcontainer
	request.URL.Path = fmt.Sprintf("/containers/%s/logs?since=0&stderr=1&stdout=1&tail=100&timestamps=0", createResponse.Id)
	request.Body = nil
	request.ContentLength = 0
	request.Method = "GET"
	request.Header.Set("Content-Type", "application/json")

	// TODO remove

	return nil
}
