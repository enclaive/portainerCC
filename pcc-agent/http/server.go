package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/portainer/agent"
	"github.com/portainer/agent/crypto"
	"github.com/portainer/agent/edge"
	"github.com/portainer/agent/exec"
	"github.com/portainer/agent/http/handler"
	"github.com/portainer/agent/kubernetes"
	httpError "github.com/portainer/libhttp/error"

	"github.com/rs/zerolog/log"
)

// APIServer is the web server exposing the API of an agent.
type APIServer struct {
	addr               string
	port               string
	systemService      agent.SystemService
	clusterService     agent.ClusterService
	signatureService   agent.DigitalSignatureService
	edgeManager        *edge.Manager
	agentTags          *agent.RuntimeConfiguration
	agentOptions       *agent.Options
	kubeClient         *kubernetes.KubeClient
	kubernetesDeployer *exec.KubernetesDeployer
	containerPlatform  agent.ContainerPlatform
	nomadConfig        agent.NomadConfig
}

// APIServerConfig represents a server configuration
// used to create a new API server
type APIServerConfig struct {
	Addr                 string
	Port                 string
	SystemService        agent.SystemService
	ClusterService       agent.ClusterService
	SignatureService     agent.DigitalSignatureService
	EdgeManager          *edge.Manager
	KubeClient           *kubernetes.KubeClient
	KubernetesDeployer   *exec.KubernetesDeployer
	RuntimeConfiguration *agent.RuntimeConfiguration
	AgentOptions         *agent.Options
	ContainerPlatform    agent.ContainerPlatform
	NomadConfig          agent.NomadConfig
}

// NewAPIServer returns a pointer to a APIServer.
func NewAPIServer(config *APIServerConfig) *APIServer {
	return &APIServer{
		addr:               config.Addr,
		port:               config.Port,
		systemService:      config.SystemService,
		clusterService:     config.ClusterService,
		signatureService:   config.SignatureService,
		edgeManager:        config.EdgeManager,
		agentTags:          config.RuntimeConfiguration,
		agentOptions:       config.AgentOptions,
		kubeClient:         config.KubeClient,
		kubernetesDeployer: config.KubernetesDeployer,
		containerPlatform:  config.ContainerPlatform,
		nomadConfig:        config.NomadConfig,
	}
}

// Start starts a new web server by listening on the specified listenAddr.
func (server *APIServer) Start(edgeMode bool) error {

	config := &handler.Config{
		SystemService:        server.systemService,
		ClusterService:       server.clusterService,
		SignatureService:     server.signatureService,
		RuntimeConfiguration: server.agentTags,
		EdgeManager:          server.edgeManager,
		KubeClient:           server.kubeClient,
		KubernetesDeployer:   server.kubernetesDeployer,
		UseTLS:               !edgeMode,
		ContainerPlatform:    server.containerPlatform,
		NomadConfig:          server.nomadConfig,
	}

	httpHandler := handler.NewHandler(config)
	httpServer := &http.Server{
		// Addr:         server.addr + ":" + server.port,
		Addr:         server.addr + ":1337",
		Handler:      httpHandler,
		ReadTimeout:  1000 * time.Second,
		WriteTimeout: 30 * time.Minute,
	}

	log.Info().
		Str("server_addr", server.addr).
		Str("server_port", server.port).
		Bool("use_tls", config.UseTLS).
		Str("api_version", agent.Version).
		Msg("starting Agent API server")

	if edgeMode {
		httpServer.Handler = server.edgeHandler(httpHandler)
		// return httpServer.ListenAndServe()
		go httpServer.ListenAndServe()
		return startTCPListener(server)
	}

	httpServer.TLSConfig = &tls.Config{
		MinVersion:   tls.VersionTLS12,
		CipherSuites: crypto.TLS12CipherSuites,
	}

	go server.securityShutdown(httpServer)
	log.Info().Msg("HIER: " + httpServer.Addr)
	// return httpServer.ListenAndServeTLS(agent.TLSCertPath, agent.TLSKeyPath)
	go httpServer.ListenAndServeTLS(agent.TLSCertPath, agent.TLSKeyPath)
	log.Info().Msg("HIER: " + httpServer.Addr)
	return startTCPListener(server)
}

// portainerCC proxy
// func proxy(conn net.Conn) {
// 	defer conn.Close()

// 	// upstream, err := net.Dial("tcp", "172.17.0.3:4444")
// 	upstream, err := net.Dial("tcp", "google.com:443")
// 	if err != nil {
// 		golog.Fatal(err)
// 		return
// 	}

// 	defer upstream.Close()

// 	go io.Copy(upstream, conn)
// 	io.Copy(conn, upstream)
// }

func startTCPListener(server *APIServer) error {
	// portainerCC proxy coordinator
	l, err := net.Listen("tcp", server.addr+":"+server.port)
	if err != nil {
		log.Print(err)
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			return err
		}
		go handleConnection(conn, server)
	}
}

func peekClientHello(reader io.Reader) (*tls.ClientHelloInfo, io.Reader, error) {
	peekedBytes := new(bytes.Buffer)
	hello, err := readClientHello(io.TeeReader(reader, peekedBytes))
	if err != nil {
		return nil, nil, err
	}
	return hello, io.MultiReader(peekedBytes, reader), nil
}

type readOnlyConn struct {
	reader io.Reader
}

func handleConnection(clientConn net.Conn, server *APIServer) {
	defer clientConn.Close()

	if err := clientConn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		log.Print(err)
		return
	}

	clientHello, clientReader, err := peekClientHello(clientConn)
	if err != nil {
		log.Print(err)
		return
	}

	if err := clientConn.SetReadDeadline(time.Time{}); err != nil {
		log.Print(err)
		return
	}

	fmt.Println("KACKMIST")
	log.Print(clientHello.ServerName)
	var target string

	if clientHello.ServerName == "coordinator" {
		// TODO hardcoded
		log.Info().Msg("CONNECTION PROXY TO COORDINATOR")
		target = "172.20.0.20:4433"

	} else {
		log.Info().Msg("DEFAULT CONNECTION / SNI")
		target = server.addr + ":1337"
	}

	backendConn, err := net.DialTimeout("tcp", target, 5*time.Second)
	if err != nil {
		log.Print(err)
		log.Err(err)
		fmt.Println("ahja")
		fmt.Println(err)
		return
	}
	defer backendConn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		io.Copy(clientConn, backendConn)
		clientConn.(*net.TCPConn).CloseWrite()
		wg.Done()
	}()
	go func() {
		io.Copy(backendConn, clientReader)
		backendConn.(*net.TCPConn).CloseWrite()
		wg.Done()
	}()

	wg.Wait()

}

func (conn readOnlyConn) Read(p []byte) (int, error)         { return conn.reader.Read(p) }
func (conn readOnlyConn) Write(p []byte) (int, error)        { return 0, io.ErrClosedPipe }
func (conn readOnlyConn) Close() error                       { return nil }
func (conn readOnlyConn) LocalAddr() net.Addr                { return nil }
func (conn readOnlyConn) RemoteAddr() net.Addr               { return nil }
func (conn readOnlyConn) SetDeadline(t time.Time) error      { return nil }
func (conn readOnlyConn) SetReadDeadline(t time.Time) error  { return nil }
func (conn readOnlyConn) SetWriteDeadline(t time.Time) error { return nil }

func readClientHello(reader io.Reader) (*tls.ClientHelloInfo, error) {
	var hello *tls.ClientHelloInfo

	err := tls.Server(readOnlyConn{reader: reader}, &tls.Config{
		GetConfigForClient: func(argHello *tls.ClientHelloInfo) (*tls.Config, error) {
			hello = new(tls.ClientHelloInfo)
			*hello = *argHello
			return nil, nil
		},
	}).Handshake()

	if hello == nil {
		return nil, err
	}

	return hello, nil
}

// !portainerCC proxy

func (server *APIServer) securityShutdown(httpServer *http.Server) {
	time.Sleep(server.agentOptions.AgentSecurityShutdown)

	if server.signatureService.IsAssociated() {
		return
	}

	log.Info().
		Dur("timeout", server.agentOptions.AgentSecurityShutdown).
		Msg("shutting down API server as no client was associated after the timeout, keeping alive to prevent restart by docker/kubernetes")

	err := httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed shutting down server")
	}
}

func (server *APIServer) edgeHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !server.edgeManager.IsKeySet() {
			httpError.WriteError(w, http.StatusForbidden, "Unable to use the unsecured agent API without Edge key", errors.New("edge key not set"))
			return
		}

		server.edgeManager.ResetActivityTimer()

		next.ServeHTTP(w, r)
	})
}
