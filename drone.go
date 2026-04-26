package gomavlinkdroneapi

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
)

type Drone struct {
	node *gomavlib.Node
}

// Connect create a node which communicates with a serial endpoint.
func (s *Drone) ConnectSerial(serialDevice string, baud int, outSystemID byte) error {
	// --- create a node which communicates with a serial endpoint.---
	// serialDevice: device serial address -> "/dev/ttyUSB0"
	// baud: baud rated of serial device -> 57600
	// outSystemID: output system ID -> 10
	s.node = &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{
				// Device: "/dev/ttyUSB0",
				// Baud:   57600,
				Device: serialDevice,
				Baud:   baud,
			},
		},
		Dialect:    ardupilotmega.Dialect,
		OutVersion: gomavlib.V2,
		// OutSystemID: 10,
		OutSystemID: outSystemID,
	}
	err := s.node.Initialize()
	return err
}

// ConnectUDPServer create a node which communicates with a UDP endpoint in server mode.
func (s *Drone) ConnectUDPServer(serverAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a UDP endpoint in server mode---
	// serverAddress: server address -> ":5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	s.node = &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			// gomavlib.EndpointUDPServer{Address: ":5600"},
			gomavlib.EndpointUDPServer{Address: serverAddress},
		},
		Dialect: ardupilotmega.Dialect,
		// OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
		OutVersion: outVersion, // change to V1 if you're unable to communicate with the target
		// OutSystemID: 10,
		OutSystemID: outSystemID,
	}
	err := s.node.Initialize()
	return err
}

// ConnectUDPServer create a node which communicates with a UDP endpoint in client mode
func (s *Drone) ConnectUDPClient(clientAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a UDP endpoint in client mode---
	// clientAddress: client address -> "1.2.3.4:5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	s.node = &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			// gomavlib.EndpointUDPClient{Address: "1.2.3.4:5600"},
			gomavlib.EndpointUDPClient{Address: clientAddress},
		},
		Dialect: ardupilotmega.Dialect,
		// OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
		OutVersion: outVersion, // change to V1 if you're unable to communicate with the target
		// OutSystemID: 10,
		OutSystemID: outSystemID,
	}
	err := s.node.Initialize()
	return err
}

// ConnectUDPServer create a node which communicates with a UDP endpoint in broadcast mode.
func (s *Drone) ConnectUDPBroadcast(broadcastAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a UDP endpoint in broadcast mode---
	// broadcastAddress: broadcast address -> "192.168.7.255:5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	s.node = &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			// gomavlib.EndpointUDPBroadcast{BroadcastAddress: "192.168.7.255:5600"},
			gomavlib.EndpointUDPBroadcast{BroadcastAddress: broadcastAddress},
		},
		Dialect: ardupilotmega.Dialect,
		// OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
		OutVersion: outVersion, // change to V1 if you're unable to communicate with the target
		// OutSystemID: 10,
		OutSystemID: outSystemID,
	}
	err := s.node.Initialize()
	return err
}

// ConnectUDPServer create a node which communicates with a TCP endpoint in server mode
func (s *Drone) ConnectTCPServer(serverAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a TCP endpoint in server mode---
	// serverAddress: serverAddress address -> ":5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	s.node = &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			// gomavlib.EndpointTCPServer{Address: ":5600"},
			gomavlib.EndpointTCPServer{Address: serverAddress},
		},
		Dialect: ardupilotmega.Dialect,
		// OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
		OutVersion: outVersion, // change to V1 if you're unable to communicate with the target
		// OutSystemID: 10,
		OutSystemID: outSystemID,
	}
	err := s.node.Initialize()
	return err
}

// ConnectUDPServer create a node which communicates with a TCP endpoint in client mode
func (s *Drone) ConnectTCPClient(clientAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a TCP endpoint in client mode---
	// clientAddress: client address -> "1.2.3.4:5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	s.node = &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			// gomavlib.EndpointTCPClient{Address: "1.2.3.4:5600"},
			gomavlib.EndpointTCPClient{Address: clientAddress},
		},
		Dialect: ardupilotmega.Dialect,
		// OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
		OutVersion: outVersion, // change to V1 if you're unable to communicate with the target
		// OutSystemID: 10,
		OutSystemID: outSystemID,
	}
	err := s.node.Initialize()
	return err
}

// ConnectUDPServer create a node which communicates with a custom TCP/TLS endpoint in client mode.
func (s *Drone) ConnectCustomClient(clientAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a custom TCP/TLS endpoint in client mode.---
	// clientAddress: client address -> "127.0.0.1:5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	s.node = &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointCustomClient{
				Connect: func(ctx context.Context) (net.Conn, error) {
					tlsConfig := &tls.Config{
						// skip checking the certificate against a CA (just set to true for simplicity for now)
						InsecureSkipVerify: true,
					}

					// return (&tls.Dialer{Config: tlsConfig}).DialContext(ctx, "tcp", "127.0.0.1:5600")
					return (&tls.Dialer{Config: tlsConfig}).DialContext(ctx, "tcp", clientAddress)
				},
				Label: "TCP/TLS:" + clientAddress,
			},
		},
		Dialect: ardupilotmega.Dialect,
		// OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
		OutVersion: outVersion, // change to V1 if you're unable to communicate with the target
		// OutSystemID: 10,
		OutSystemID: outSystemID,
	}
	err := s.node.Initialize()
	return err
}

// ConnectUDPServer create a node which communicates with a custom TCP/TLS endpoint in server mode.
func (s *Drone) ConnectCustomServer(listenAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a custom TCP/TLS endpoint in server mode.---
	// listenAddress: listenAddress -> ":5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	var rtn error
	errcert := ensureCertsExist()
	if errcert == nil {
		node := &gomavlib.Node{
			Endpoints: []gomavlib.EndpointConf{
				gomavlib.EndpointCustomServer{
					Listen: func() (net.Listener, error) {
						// Loads the certificate and key from the generated certs dir
						cert, err := tls.LoadX509KeyPair("certs/cert.pem", "certs/key.pem")
						if err != nil {
							return nil, err
						}

						// return tls.Listen("tcp", ":5600", &tls.Config{
						// 	Certificates: []tls.Certificate{cert},
						// })
						return tls.Listen("tcp", listenAddress, &tls.Config{
							Certificates: []tls.Certificate{cert},
						})
					},
					Label: "TCP/TLS",
				},
			},
			Dialect: ardupilotmega.Dialect,
			// OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
			OutVersion: outVersion, // change to V1 if you're unable to communicate with the target
			// OutSystemID: 10,
			OutSystemID: outSystemID,
		}
		rtn = node.Initialize()
	} else {
		rtn = errcert
	}

	return rtn
}

func (s *Drone) Close() {
	s.node.Close()
}

func ensureCertsExist() error {
	// Check if cert.pem exists
	if _, err := os.Stat("certs/cert.pem"); os.IsNotExist(err) {
		fmt.Println("cert.pem not found. Generating certificates...")
		return generateCertAndKey()
	}

	// Check if key.pem exists
	if _, err := os.Stat("certs/key.pem"); os.IsNotExist(err) {
		fmt.Println("key.pem not found. Generating certificates...")
		return generateCertAndKey()
	}

	return nil
}

// GenerateCertAndKey generates a self-signed certificate and private key, saving them to the certs/ directory.
func generateCertAndKey() error {
	// Create the certs directory if it doesn't exist
	err := os.MkdirAll("certs", os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create certs directory: %w", err)
	}

	// Generate RSA private key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "gomavlib",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour), // valid for 1 year
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		BasicConstraintsValid: true,
	}

	// Create the certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	// Save the certificate
	certOut, err := os.Create("certs/cert.pem")
	if err != nil {
		return fmt.Errorf("failed to create cert.pem: %w", err)
	}
	defer certOut.Close()

	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	if err != nil {
		return fmt.Errorf("failed to encode certificate to PEM: %w", err)
	}

	// Save the private key
	keyOut, err := os.Create("certs/key.pem")
	if err != nil {
		return fmt.Errorf("failed to create key.pem: %w", err)
	}
	defer keyOut.Close()

	err = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	if err != nil {
		return fmt.Errorf("failed to encode private key to PEM: %w", err)
	}

	fmt.Println("cert.pem and key.pem generated in the 'certs/' directory.")
	return nil
}
