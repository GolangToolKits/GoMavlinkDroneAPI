package gomavlinkdroneapi

import (
	"context"
	"errors"
	"log"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
)

type DroneAPI struct {
	drone *Drone
}

// ConnectSerial  ConnectSerial
func (s *DroneAPI) ConnectSerial(serialDevice string, baud int, outSystemID byte) error {
	// --- create a node which communicates with a serial endpoint.---
	// serialDevice: device serial address -> "/dev/ttyUSB0"
	// baud: baud rated of serial device -> 57600
	// outSystemID: output system ID -> 10

	s.drone = &Drone{}
	return s.drone.ConnectSerial(serialDevice, baud, outSystemID)
}

// ConnectUDPServer  ConnectUDPServer
func (s *DroneAPI) ConnectUDPServer(serverAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a UDP endpoint in server mode---
	// serverAddress: server address -> ":5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	s.drone = &Drone{}
	return s.drone.ConnectUDPServer(serverAddress, outVersion, outSystemID)
}

// ConnectUDPClient ConnectUDPClient
func (s *DroneAPI) ConnectUDPClient(clientAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a UDP endpoint in client mode---
	// clientAddress: client address -> "1.2.3.4:5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	s.drone = &Drone{}
	return s.drone.ConnectUDPClient(clientAddress, outVersion, outSystemID)
}

// ConnectUDPBroadcast ConnectUDPBroadcast
func (s *DroneAPI) ConnectUDPBroadcast(broadcastAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a UDP endpoint in broadcast mode---
	// broadcastAddress: broadcast address -> "192.168.7.255:5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	s.drone = &Drone{}
	return s.drone.ConnectUDPBroadcast(broadcastAddress, outVersion, outSystemID)
}

// ConnectTCPServer ConnectTCPServer
func (s *DroneAPI) ConnectTCPServer(serverAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a TCP endpoint in server mode---
	// serverAddress: serverAddress address -> ":5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	s.drone = &Drone{}
	return s.drone.ConnectTCPServer(serverAddress, outVersion, outSystemID)
}

// ConnectTCPClient ConnectTCPClient
func (s *DroneAPI) ConnectTCPClient(clientAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a TCP endpoint in client mode---
	// clientAddress: client address -> "1.2.3.4:5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	s.drone = &Drone{}
	return s.drone.ConnectTCPClient(clientAddress, outVersion, outSystemID)
}

// ConnectCustomClient ConnectCustomClient
func (s *DroneAPI) ConnectCustomClient(clientAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a custom TCP/TLS endpoint in client mode.---
	// clientAddress: client address -> "127.0.0.1:5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	s.drone = &Drone{}
	return s.drone.ConnectCustomClient(clientAddress, outVersion, outSystemID)
}

// ConnectCustomServer ConnectCustomServer
func (s *DroneAPI) ConnectCustomServer(listenAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a custom TCP/TLS endpoint in server mode.---
	// listenAddress: listenAddress -> ":5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	s.drone = &Drone{}
	return s.drone.ConnectCustomServer(listenAddress, outVersion, outSystemID)
}

// func (s *DroneAPI) Land() (*action.LandResponse, error) {
// 	// lr, err := s.Drone.action.Land(context.Background())
// 	// return lr, err
// 	return nil, nil
// }

// func (s *DroneAPI) ConnectionState() (<-chan *core.ConnectionState, error) {
// 	// cs, err := s.Drone.core.ConnectionState(context.Background())
// 	// return cs, err
// 	return nil, nil
// }

// func (s *DroneAPI) UploadGeofence(geofenceData *geofence.GeofenceData) (*geofence.UploadGeofenceResponse, error) {
// 	// gfr, err := s.Drone.geofence.UploadGeofence(context.Background(), geofenceData)
// 	// return gfr, err
// 	return nil, nil
// }

// IsDroneConnected checks to see if drone is connected
func (s *DroneAPI) IsDroneConnected(ctx context.Context) (bool, error) {

	log.Println("Is Drone connected, Waiting for drone heartbeat...")

	for {
		select {
		// 1. Listen for the external cancellation or timeout signal
		case <-ctx.Done():
			log.Println("Connection attempt aborted by external signal.")
			return false, ctx.Err()

		// 2. Safely read from the event channel
		case evt, ok := <-s.drone.node.Events():
			if !ok {
				log.Println("Connection failed: Event channel closed unexpectedly.")
				s.drone.node.Close()
				return false, errors.New("event channel closed unexpectedly")
			}

			// Ensure the event is a frame
			frm, ok := evt.(*gomavlib.EventFrame)
			if !ok {
				continue // Skip other event types like EventChannelOpen
			}

			// Type assertion check for Heartbeat instead of ID == 0
			switch msg := frm.Message().(type) {
			case *ardupilotmega.MessageHeartbeat:
				log.Printf("🎉 Success! Drone connected. System ID: %d, Autopilot: %d", frm.SystemID(), msg.Autopilot)
				return true, nil
			}
		}
	}
}

// ListenToDroneEvents gets the event for listing in the calling app
func (s *DroneAPI) ListenToDroneEvents() chan gomavlib.Event {
	// Once initialized, you must loop over the node.Events() channel.
	// This handles incoming messages and maintains the internal buffer clearing
	return s.drone.node.Events()
}

// Close ends connection to vehicle
func (s *DroneAPI) Close() {
	s.drone.Close()
}

// New ghts a new api reference
func (s *DroneAPI) New() API {
	return s
}

// Once initialized, you must loop over the node.Events() channel.
// This handles incoming messages and maintains the internal buffer clearing

// func listenToDrone(node *gomavlib.Node) {
//     for evt := range node.Events() {
//         switch e := evt.(type) {

//         // Triggered when a message is successfully decoded
//         case *gomavlib.EventFrame:
//             msg := e.Message()

//             // Print the raw message ID and its contents
//             log.Printf("Received msg from System %d: ID %d (%+v)",
//                 e.SystemID(), msg.GetID(), msg)

//         // Triggered when a new physical connection/node is established
//         case *gomavlib.EventChannelOpen:
//             log.Printf("Channel opened: %s", e.Channel)

//         // Triggered when a drone disconnects
//         case *gomavlib.EventChannelClose:
//             log.Printf("Channel closed: %s", e.Channel)
//         }
//     }
// }

// ⚡ Bonus Tip for ArduPilot: If you are connecting to an ArduPilot device, it usually requires a stream request before it floods you with data. You can automate this by adding StreamRequests: true inside your gomavlib.NodeConf.
// Go Packages
// Go Packages
//  +2
// What flight controller (PX4 or ArduPilot) and physical connection are you using for this project?

// ⚠️ Important TCP MAVLink Considerations
// Reconnection Handling: The gomavlib library features automated background reconnection mechanics in case the drone briefly drops its Wi-Fi/cellular connection.
// Latency vs Reliability: As noted by expert forums on GitHub, TCP demands heavy data acknowledgement and forces data re-transmissions if signals weaken. While this prevents data loss, it can cause control lag compared to UDP.
// Port Norms: Custom setups vary, but typical MAVLink simulation defaults often rely on ports like 5760 or 5762.
// GitHub
// GitHub
//  +4
// Are you connecting to a physical drone network or operating a local software-in-the-loop (SITL) simulation?

// In Go, calling gomavlib.NewNode() only checks that your local network adapter or serial port is accessible. It does not verify that a physical drone is actively listening, powered on, or communicating back to you.
// MAVSDK
// MAVSDK
//  +1
// To confirm a successful connection, you must wait to receive a HEARTBEAT message from the drone's System ID

// 💓 The MAVLink Heartbeat Rule
// All MAVLink-compatible flight controllers (like PX4 and ArduPilot) broadcast a HEARTBEAT message once every second. If you aren't receiving heartbeats, you are not connected.
// Blue Robotics
// Blue Robotics
//  +4
// Here is how you can perform a connection check in gomavlib with a timeout fallback:
