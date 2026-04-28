package gomavlinkdroneapi_test

import (
	"testing"

	gomavlinkdroneapi "github.com/GolangToolKits/GoMavlinkDroneAPI"
	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

func TestDroneAPI_UploadMission(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		missionItems    *[]ardupilotmega.MessageMissionItemInt
		targetSystem    uint8
		targetComponent uint8
		want            bool
		clientAddress   string
		outVersion      gomavlib.Version
		outSystemID     byte
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			name:          "test 1",
			clientAddress: "1.2.3.4:5600",
			outVersion:    gomavlib.V2,
			outSystemID:   10,
			wantErr:       false,
			missionItems: &[]ardupilotmega.MessageMissionItemInt{
				{
					TargetSystem:    1,
					TargetComponent: 1,
					Seq:             0,
					Frame:           common.MAV_FRAME_GLOBAL_RELATIVE_ALT,
					Command:         common.MAV_CMD_NAV_TAKEOFF,
					Current:         1,
					Autocontinue:    1,
					Z:               10, // Take off to 10 meters
				},
				{
					TargetSystem:    1,
					TargetComponent: 1,
					Seq:             1,
					Frame:           common.MAV_FRAME_GLOBAL_RELATIVE_ALT,
					Command:         common.MAV_CMD_NAV_WAYPOINT,
					Current:         0,
					Autocontinue:    1,
					X:               473977418, // Latitude * 1e7
					Y:               85443833,  // Longitude * 1e7
					Z:               15,        // Altitude (meters)
				},
			},
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
			got := s.UploadMission(tt.missionItems, tt.targetSystem, tt.targetComponent)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("UploadMission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDroneAPI_DownloadMissions(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		targetSystem    uint8
		targetComponent uint8
		want            map[uint16]*ardupilotmega.MessageMissionItemInt
		wantErr         bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var s gomavlinkdroneapi.DroneAPI
			got, gotErr := s.DownloadMissions(tt.targetSystem, tt.targetComponent)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("DownloadMissions() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("DownloadMissions() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("DownloadMissions() = %v, want %v", got, tt.want)
			}
		})
	}
}
