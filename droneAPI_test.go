package gomavlinkdroneapi_test

import (
	"testing"

	gomavlinkdroneapi "github.com/GolangToolKits/GoMavlinkDroneAPI"
	"github.com/bluenviron/gomavlib/v3"
)

func TestDroneAPI_ConnectSerial(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		serialDevice  string
		baud          int
		outSystemID   byte
		wantErr       bool
		testHeartBeat bool
	}{
		// TODO: Add test cases.
		{
			name:         "test 1",
			serialDevice: "/dev/ttyUSB0",
			baud:         57600,
			outSystemID:  10,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var s gomavlinkdroneapi.DroneAPI
			gotErr := s.ConnectSerial(tt.serialDevice, tt.baud, tt.outSystemID)
			if tt.testHeartBeat {
				con := s.IsDroneConnected()
				if !con {
					t.Fatal("ConnectSerial() succeeded not connected")
				}
			}

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ConnectSerial() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ConnectSerial() succeeded unexpectedly")
			}
		})
	}
}

func TestDroneAPI_ConnectUDPServer(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		serverAddress string
		outVersion    gomavlib.Version
		outSystemID   byte
		wantErr       bool
	}{
		// TODO: Add test cases.
		{
			name:          "test 1",
			serverAddress: ":5600",
			outVersion:    gomavlib.V2,
			outSystemID:   10,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var s gomavlinkdroneapi.DroneAPI
			gotErr := s.ConnectUDPServer(tt.serverAddress, tt.outVersion, tt.outSystemID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ConnectUDPServer() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ConnectUDPServer() succeeded unexpectedly")
			}
		})
	}
}

func TestDroneAPI_ConnectUDPClient(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		clientAddress string
		outVersion    gomavlib.Version
		outSystemID   byte
		wantErr       bool
	}{
		// TODO: Add test cases.
		{
			name:          "test 1",
			clientAddress: "1.2.3.4:5600",
			outVersion:    gomavlib.V2,
			outSystemID:   10,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var s gomavlinkdroneapi.DroneAPI
			gotErr := s.ConnectUDPClient(tt.clientAddress, tt.outVersion, tt.outSystemID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ConnectUDPClient() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ConnectUDPClient() succeeded unexpectedly")
			}
		})
	}
}

func TestDroneAPI_ConnectUDPBroadcast(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		broadcastAddress string
		outVersion       gomavlib.Version
		outSystemID      byte
		wantErr          bool
	}{
		// TODO: Add test cases.
		{
			name:             "test 1",
			broadcastAddress: "192.168.7.255:5600",
			outVersion:       gomavlib.V2,
			outSystemID:      10,
			wantErr:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var s gomavlinkdroneapi.DroneAPI
			gotErr := s.ConnectUDPBroadcast(tt.broadcastAddress, tt.outVersion, tt.outSystemID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ConnectUDPBroadcast() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ConnectUDPBroadcast() succeeded unexpectedly")
			}
		})
	}
}

func TestDroneAPI_ConnectTCPServer(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		serverAddress string
		outVersion    gomavlib.Version
		outSystemID   byte
		wantErr       bool
	}{
		// TODO: Add test cases.
		{
			name:          "test 1",
			serverAddress: ":5600",
			outVersion:    gomavlib.V2,
			outSystemID:   10,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var s gomavlinkdroneapi.DroneAPI
			gotErr := s.ConnectTCPServer(tt.serverAddress, tt.outVersion, tt.outSystemID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ConnectTCPServer() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ConnectTCPServer() succeeded unexpectedly")
			}
		})
	}
}

func TestDroneAPI_ConnectTCPClient(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		clientAddress string
		outVersion    gomavlib.Version
		outSystemID   byte
		wantErr       bool
	}{
		// TODO: Add test cases.
		{
			name:          "test 1",
			clientAddress: "1.2.3.4:5600",
			outVersion:    gomavlib.V2,
			outSystemID:   10,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var s gomavlinkdroneapi.DroneAPI
			gotErr := s.ConnectTCPClient(tt.clientAddress, tt.outVersion, tt.outSystemID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ConnectTCPClient() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ConnectTCPClient() succeeded unexpectedly")
			}
		})
	}
}

func TestDroneAPI_ConnectCustomClient(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		clientAddress string
		outVersion    gomavlib.Version
		outSystemID   byte
		wantErr       bool
	}{
		// TODO: Add test cases.
		{
			name:          "test 1",
			clientAddress: "1.2.3.4:5600",
			outVersion:    gomavlib.V2,
			outSystemID:   10,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var s gomavlinkdroneapi.DroneAPI
			gotErr := s.ConnectCustomClient(tt.clientAddress, tt.outVersion, tt.outSystemID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ConnectCustomClient() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ConnectCustomClient() succeeded unexpectedly")
			}
		})
	}
}

func TestDroneAPI_ConnectCustomServer(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		listenAddress string
		outVersion    gomavlib.Version
		outSystemID   byte
		wantErr       bool
	}{
		// TODO: Add test cases.
		{
			name:          "test 1",
			listenAddress: ":5600",
			outVersion:    gomavlib.V2,
			outSystemID:   10,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var s gomavlinkdroneapi.DroneAPI
			gotErr := s.ConnectCustomServer(tt.listenAddress, tt.outVersion, tt.outSystemID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ConnectCustomServer() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ConnectCustomServer() succeeded unexpectedly")
			}
		})
	}
}
