package gomavlinkdroneapi_test

import (
	"testing"

	gomavlinkdroneapi "github.com/GolangToolKits/GoMavlinkDroneAPI"
	"github.com/bluenviron/gomavlib/v3"
)

func TestDroneAPI_ArmDisarm(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		//imputs for UDPClient node
		clientAddress string
		outVersion    gomavlib.Version
		outSystemID   byte
		// Named input parameters for target function.
		armCommand      float32
		disarmCommand   float32
		targetSystem    uint8
		targetComponent uint8
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			name:          "test 1",
			clientAddress: "1.2.3.4:5600",
			outVersion:    gomavlib.V2,
			outSystemID:   10,
			// commands
			armCommand:      1,
			disarmCommand:   0,
			targetSystem:    1,
			targetComponent: 1,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var ss gomavlinkdroneapi.DroneAPI
			s := ss.New()
			gotContErr := s.ConnectUDPClient(tt.clientAddress, tt.outVersion, tt.outSystemID)
			if gotContErr != nil {
				if !tt.wantErr {
					t.Errorf("ConnectUDPClient() failed: %v", gotContErr)
				}
				return
			}
			armErr := s.ArmDisarm(tt.armCommand, tt.targetSystem, tt.targetComponent)
			if armErr != nil {
				if !tt.wantErr {
					t.Errorf("ArmDisarm() failed: %v", armErr)
				}
				return
			}
			disarmErr := s.ArmDisarm(tt.disarmCommand, tt.targetSystem, tt.targetComponent)
			if disarmErr != nil {
				if !tt.wantErr {
					t.Errorf("ArmDisarm() failed: %v", disarmErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ArmDisarm() succeeded unexpectedly")
			}
		})
	}
}
