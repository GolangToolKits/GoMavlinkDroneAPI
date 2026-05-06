package gomavlinkdroneapi_test

import (
	"context"
	"fmt"
	"testing"

	gomavlinkdroneapi "github.com/GolangToolKits/GoMavlinkDroneAPI"
	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

func TestDroneAPI_ClearGeofence(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		targetSystem    uint8
		targetComponent uint8
		want            bool
		clientAddress   string
		outVersion      gomavlib.Version
		outSystemID     byte
		wantErr         bool
	}{
		// TODO: Add test cases.
		// {
		// 	name:          "test 1",
		// 	clientAddress: "1.2.3.4:5600",
		// 	outVersion:    gomavlib.V2,
		// 	outSystemID:   10,
		// 	wantErr:       false,
		// },
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
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			got := s.ClearGeofence(ctx, tt.targetSystem, tt.targetComponent)
			// TODO: update the condition below to compare got with tt.want.
			if got != false {
				t.Errorf("ClearGeofence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDroneAPI_UploadGeofence(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		// newFence        []gomavlinkdroneapi.GeoFence
		geofence        []ardupilotmega.MessageMissionItemInt
		targetSystem    uint8
		targetComponent uint8
		want            bool
		wantErr         bool
		clientAddress   string
		outVersion      gomavlib.Version
		outSystemID     byte
	}{
		// TODO: Add test cases.
		{
			name:            "test 1",
			clientAddress:   "127.0.0.1:5760",
			outVersion:      gomavlib.V2,
			outSystemID:     255,
			wantErr:         false,
			targetSystem:    1,
			targetComponent: 1,
			geofence: []ardupilotmega.MessageMissionItemInt{
				{
					Seq:         0,
					Frame:       common.MAV_FRAME_GLOBAL,
					Command:     common.MAV_CMD_NAV_FENCE_POLYGON_VERTEX_INCLUSION,
					Param1:      5, // Updated to 5 total points
					X:           33992000,
					Y:           -84785000,
					MissionType: common.MAV_MISSION_TYPE_FENCE,
				},
				{
					Seq:         1,
					Frame:       common.MAV_FRAME_GLOBAL,
					Command:     common.MAV_CMD_NAV_FENCE_POLYGON_VERTEX_INCLUSION,
					X:           33992000,
					Y:           -84783000,
					MissionType: common.MAV_MISSION_TYPE_FENCE,
				},
				{
					Seq:         2,
					Frame:       common.MAV_FRAME_GLOBAL,
					Command:     common.MAV_CMD_NAV_FENCE_POLYGON_VERTEX_INCLUSION,
					X:           33990000,
					Y:           -84783000,
					MissionType: common.MAV_MISSION_TYPE_FENCE,
				},
				{
					Seq:         3,
					Frame:       common.MAV_FRAME_GLOBAL,
					Command:     common.MAV_CMD_NAV_FENCE_POLYGON_VERTEX_INCLUSION,
					X:           33990000,
					Y:           -84785000,
					MissionType: common.MAV_MISSION_TYPE_FENCE,
				},
				{
					// 5th POINT: Must match Seq 0 exactly to close the polygon
					Seq:         4,
					Frame:       common.MAV_FRAME_GLOBAL,
					Command:     common.MAV_CMD_NAV_FENCE_POLYGON_VERTEX_INCLUSION,
					X:           33992000,  // Same as Seq 0
					Y:           -84785000, // Same as Seq 0
					MissionType: common.MAV_MISSION_TYPE_FENCE,
				},
			},
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
			con, err := s.IsDroneConnected(ctx)
			if !con {
				fmt.Print(err)
				t.Fatal("ConnectSerial() succeeded not connected")
			}

			modeErr := s.SetMode(tt.targetSystem, tt.targetComponent, gomavlinkdroneapi.MODE_CUSTOM, uint32(gomavlinkdroneapi.MODE_GUIDED))
			if modeErr != nil {
				if !tt.wantErr {
					t.Errorf("SetMode() failed: %v", modeErr)
				}
				return
			}
			//ctx, cancel := context.WithCancel(context.Background())
			//defer cancel()
			got, gotErr := s.UploadGeofence(ctx, tt.geofence, tt.targetSystem, tt.targetComponent)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("UploadGeofence() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("UploadGeofence() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if got {
				t.Errorf("UploadGeofence() = %v, want %v", got, tt.want)
			}
			s.Close()
		})
	}
}
