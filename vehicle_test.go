package gomavlinkdroneapi_test

import (
	"context"
	"testing"

	gomavlinkdroneapi "github.com/GolangToolKits/GoMavlinkDroneAPI"
	"github.com/bluenviron/gomavlib/v3"
)

func TestDroneAPI_GetConnectedVehicle(t *testing.T) {
	tests := []struct {
		name            string // description of this test case
		want            byte
		want2           byte
		targetSystem    uint8
		targetComponent uint8
		clientAddress   string
		outVersion      gomavlib.Version
		outSystemID     byte
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			name:          "test 1",
			clientAddress: "127.0.0.1:5760",
			outVersion:    gomavlib.V2,
			outSystemID:   255,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var ss gomavlinkdroneapi.DroneAPI
			s := ss.New()
			gotErr := s.ConnectTCPClient(tt.clientAddress, tt.outVersion, tt.outSystemID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ConnectTCPClient() failed: %v", gotErr)
				}
				return
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			con, conErr := s.IsDroneConnected(ctx)
			if conErr != nil {
				if !tt.wantErr {
					t.Errorf("IsDroneConnected() failed: %v", conErr)
				}
				return
			}
			if con {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				got, got2 := s.GetConnectedVehicle(ctx)
				// TODO: update the condition below to compare got with tt.want.
				if got == 0 {
					t.Errorf("GetConnectedVehicle() = %v, want %v", got, tt.want)
				}
				if got == 0 {
					t.Errorf("GetConnectedVehicle() = %v, want %v", got2, tt.want2)
				}
			}
			s.Close()

		})
	}
}
