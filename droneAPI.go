package gomavlinkdroneapi

import (
	"context"
	"log"
	"time"

	"github.com/bluenviron/gomavlib/v3"
)

type DroneAPI struct {
	drone *Drone
}

func (s *DroneAPI) ConnectSerial(serialDevice string, baud int, outSystemID byte) error {
	// --- create a node which communicates with a serial endpoint.---
	// serialDevice: device serial address -> "/dev/ttyUSB0"
	// baud: baud rated of serial device -> 57600
	// outSystemID: output system ID -> 10

	s.drone = &Drone{}
	return s.drone.ConnectSerial(serialDevice, baud, outSystemID)
}
func (s *DroneAPI) ConnectUDPServer(serverAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a UDP endpoint in server mode---
	// serverAddress: server address -> ":5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	s.drone = &Drone{}
	return s.drone.ConnectUDPServer(serverAddress, outVersion, outSystemID)
}

func (s *DroneAPI) ConnectUDPClient(clientAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a UDP endpoint in client mode---
	// clientAddress: client address -> "1.2.3.4:5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	s.drone = &Drone{}
	return s.drone.ConnectUDPClient(clientAddress, outVersion, outSystemID)
}

func (s *DroneAPI) ConnectUDPBroadcast(broadcastAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a UDP endpoint in broadcast mode---
	// broadcastAddress: broadcast address -> "192.168.7.255:5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	s.drone = &Drone{}
	return s.drone.ConnectUDPBroadcast(broadcastAddress, outVersion, outSystemID)
}

func (s *DroneAPI) ConnectTCPServer(serverAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a TCP endpoint in server mode---
	// serverAddress: serverAddress address -> ":5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	s.drone = &Drone{}
	return s.ConnectTCPServer(serverAddress, outVersion, outSystemID)
}

func (s *DroneAPI) ConnectTCPClient(clientAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a TCP endpoint in client mode---
	// clientAddress: client address -> "1.2.3.4:5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	s.drone = &Drone{}
	return s.drone.ConnectTCPClient(clientAddress, outVersion, outSystemID)
}

func (s *DroneAPI) ConnectCustomClient(clientAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a custom TCP/TLS endpoint in client mode.---
	// clientAddress: client address -> "127.0.0.1:5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	s.drone = &Drone{}
	return s.drone.ConnectCustomClient(clientAddress, outVersion, outSystemID)
}

func (s *DroneAPI) ConnectCustomServer(listenAddress string, outVersion gomavlib.Version, outSystemID byte) error {
	// ---create a node which communicates with a custom TCP/TLS endpoint in server mode.---
	// listenAddress: listenAddress -> ":5600"
	// outVersion: out version -> gomavlib.V2 or change to V1 if you're unable to communicate with the target
	// outSystemID: output system ID -> 10
	s.drone = &Drone{}
	return s.drone.ConnectCustomServer(listenAddress, outVersion, outSystemID)
}

// func (s *DroneAPI) Arm() (*action.ArmResponse, error) {
// 	// ar, err := s.Drone.action.Arm(context.Background())
// 	// return ar, err
// 	return nil, nil
// }

// func (s *DroneAPI) Takeoff() (*action.TakeoffResponse, error) {
// 	// tr, err := s.Drone.action.Takeoff(context.Background())
// 	// return tr, err
// 	return nil, nil
// }

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

func (s *DroneAPI) IsDroneConnected() bool {

	log.Println("Node started. Waiting for drone heartbeat...")

	// Create a context that times out after 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connected := false

	// Loop through events until we find a heartbeat or the timeout expires
	for {
		select {
		case <-ctx.Done():
			log.Fatal("Connection failed: No heartbeat received from drone.")
		case evt, ok := <-s.drone.node.Events():
			if !ok {
				log.Fatal("Event channel closed.")
			}

			if frm, ok := evt.(*gomavlib.EventFrame); ok {
				// Check if the message is a HEARTBEAT (Message ID 0)
				if frm.Message().GetID() == 0 {
					log.Printf("🎉 Success! Drone connected. System ID: %d", frm.SystemID())
					connected = true
					break
				}
			}
		}
		if connected {
			break
		}
	}
	return connected
}

func (s *DroneAPI) ListenToDroneEvents() chan gomavlib.Event {
	// Once initialized, you must loop over the node.Events() channel.
	// This handles incoming messages and maintains the internal buffer clearing
	return s.drone.node.Events()
}

func (s *DroneAPI) Close() {
	s.drone.Close()
}

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
