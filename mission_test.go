package gomavlinkdroneapi_test

import (
	"testing"

	gomavlinkdroneapi "github.com/GolangToolKits/GoMavlinkDroneAPI"
	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
)

func TestDroneAPI_UploadMission(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		droneChannel *gomavlib.Channel
		mission      []ardupilotmega.MessageMissionItemInt
		want         bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var s gomavlinkdroneapi.DroneAPI
			got := s.UploadMission(tt.droneChannel, tt.mission)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("UploadMission() = %v, want %v", got, tt.want)
			}
		})
	}
}
