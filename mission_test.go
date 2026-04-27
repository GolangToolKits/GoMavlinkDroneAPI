package gomavlinkdroneapi_test

import (
	"testing"

	gomavlinkdroneapi "github.com/GolangToolKits/GoMavlinkDroneAPI"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
)

func TestDroneAPI_UploadMission(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		missionItems    []ardupilotmega.MessageMissionItemInt
		targetSystem    uint8
		targetComponent uint8
		want            bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var s gomavlinkdroneapi.DroneAPI
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
