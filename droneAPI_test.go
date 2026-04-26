package gomavlinkdroneapi_test

import (
	"testing"

	gomavlinkdroneapi "github.com/GolangToolKits/GoMavlinkDroneAPI"
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
