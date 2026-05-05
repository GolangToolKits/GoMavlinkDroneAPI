package gomavlinkdroneapi_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	gomavlinkdroneapi "github.com/GolangToolKits/GoMavlinkDroneAPI"
	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

func TestDroneAPI_UploadMission(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		missionItems           []ardupilotmega.MessageMissionItemInt
		targetSystem           uint8
		targetComponent        uint8
		want                   bool
		clientAddress          string
		outVersion             gomavlib.Version
		outSystemID            byte
		wantErr                bool
		overrideMission        bool
		missionOverrideCommand *common.MessageCommandLong
		disarmCommand          bool
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
			missionItems: []ardupilotmega.MessageMissionItemInt{
				{
					TargetSystem:    1,
					TargetComponent: 1,
					Seq:             0,
					Frame:           common.MAV_FRAME_GLOBAL_RELATIVE_ALT,
					Command:         common.MAV_CMD_NAV_TAKEOFF,
					Current:         1,
					Autocontinue:    1,
					Z:               1, // Take off to 10 meters
				},
				{
					TargetSystem:    1,
					TargetComponent: 1,
					Seq:             1,
					Frame:           common.MAV_FRAME_GLOBAL_RELATIVE_ALT,
					Command:         common.MAV_CMD_NAV_WAYPOINT,
					Current:         0,
					Autocontinue:    1,
					// X:               473977418, // Latitude * 1e7
					X: 33991097,  // Latitude * 1e7
					Y: -84784404, // Longitude * 1e7
					Z: 10,        // Altitude (meters)
				},
			},
		},
		{
			name:            "test 2",
			clientAddress:   "127.0.0.1:5760",
			outVersion:      gomavlib.V2,
			outSystemID:     255,
			wantErr:         false,
			targetSystem:    1,
			targetComponent: 1,
			overrideMission: true,
			disarmCommand:   false,
			missionItems: []ardupilotmega.MessageMissionItemInt{
				{
					TargetSystem:    1,
					TargetComponent: 1,
					Seq:             0,
					Frame:           common.MAV_FRAME_GLOBAL_RELATIVE_ALT,
					Command:         common.MAV_CMD_NAV_TAKEOFF,
					Current:         1,
					Autocontinue:    1,
					Z:               1, // Take off to 10 meters
				},
				{
					TargetSystem:    1,
					TargetComponent: 1,
					Seq:             1,
					Frame:           common.MAV_FRAME_GLOBAL_RELATIVE_ALT,
					Command:         common.MAV_CMD_NAV_WAYPOINT,
					Current:         0,
					Autocontinue:    1,
					// X:               473977418, // Latitude * 1e7
					X: 33991097,  // Latitude * 1e7
					Y: -84784404, // Longitude * 1e7
					Z: 10,        // Altitude (meters)
				},
			},
			missionOverrideCommand: &common.MessageCommandLong{
				TargetSystem:    1, // System ID of the drone
				TargetComponent: 1, // Component ID of the drone
				// Command:         common.MAV_CMD_DO_REPOSITION,
				Command: common.MAV_CMD_NAV_LOITER_UNLIM,
				Param1:  float32(common.MAV_MODE_FLAG_CUSTOM_MODE_ENABLED),
				Param2:  4, // Example custom mode (HOLD in PX4)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("Running test: %s\n", tt.name)
			// TODO: construct the receiver type.
			var ss gomavlinkdroneapi.DroneAPI
			s := ss.New()
			// gotContErr := s.ConnectUDPClient(tt.clientAddress, tt.outVersion, tt.outSystemID)
			// if gotContErr != nil {
			// 	if !tt.wantErr {
			// 		t.Errorf("ConnectUDPClient() failed: %v", gotContErr)
			// 	}
			// 	return
			// }
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
			modeErr := s.SetMode(tt.targetSystem, tt.targetComponent, gomavlinkdroneapi.MODE_CUSTOM, uint32(gomavlinkdroneapi.MODE_AUTO))
			if modeErr != nil {
				if !tt.wantErr {
					t.Errorf("SetMode() failed: %v", modeErr)
				}
				return
			}
			//ctx, cancel := context.WithCancel(context.Background())
			//defer cancel()
			got := s.UploadMission(ctx, tt.missionItems, tt.targetSystem, tt.targetComponent)
			// TODO: update the condition below to compare got with tt.want.
			if got != nil {
				t.Errorf("UploadMission() = %v, want %v", got, tt.want)
			}

			gotd, dowerr := s.DownloadMissions(ctx, tt.targetSystem, tt.targetComponent)
			// TODO: update the condition below to compare got with tt.want.
			if dowerr != nil {
				t.Errorf("DownloadMissions() = %v, want %v", gotd, tt.want)
			} else {
				for i := 0; i < len(gotd); i++ {
					fmt.Printf("Downloaded mission %d: %v\n", i, gotd[i])
				}
				//fmt.Printf("Downloaded missions: %v\n", gotd[0])
			}

			gotma := s.StartMission(ctx, tt.targetSystem, tt.targetComponent)
			// TODO: update the condition below to compare got with tt.want.
			if gotma != nil {
				t.Errorf("StartMission() = %v, want %v", gotma, tt.want)
			}

			if !tt.overrideMission {
				gotmon, monerr := s.MonitorMission(ctx, 2, tt.targetSystem)
				// TODO: update the condition below to compare got with tt.want.
				if monerr != nil {
					t.Errorf("MonitorMission() = %v, want %v", gotmon, tt.want)
				}
				homeErr := s.ReturnHome(ctx, tt.targetSystem, tt.targetComponent)
				if homeErr != nil {
					if !tt.wantErr {
						t.Errorf("ReturnHome() failed: %v", homeErr)
					}
					return
				}
				time.Sleep(120 * time.Second)
				lndErr := s.Land(ctx, tt.targetSystem, tt.targetComponent)
				if lndErr != nil {
					if !tt.wantErr {
						t.Errorf("Land() failed: %v", lndErr)
					}
					return
				}
				time.Sleep(30 * time.Second)
				disarmErr := s.ArmDisarm(ctx, tt.disarmCommand, tt.targetSystem, tt.targetComponent)
				if disarmErr != nil {
					if !tt.wantErr {
						t.Errorf("ArmDisarm() failed: %v", disarmErr)
					}
					return
				}
			} else {
				orderr := s.OverrideMissionAndHover(ctx, tt.missionOverrideCommand)
				// TODO: update the condition below to compare got with tt.want.
				if orderr != nil {
					t.Errorf("OverrideMissionAndHover() = %v, want %v", orderr, tt.want)
				}
				homeErr := s.ReturnHome(ctx, tt.targetSystem, tt.targetComponent)
				if homeErr != nil {
					if !tt.wantErr {
						t.Errorf("ReturnHome() failed: %v", homeErr)
					}
					return
				}
				time.Sleep(120 * time.Second)
				lndErr := s.Land(ctx, tt.targetSystem, tt.targetComponent)
				if lndErr != nil {
					if !tt.wantErr {
						t.Errorf("Land() failed: %v", lndErr)
					}
					return
				}
				time.Sleep(30 * time.Second)
				disarmErr := s.ArmDisarm(ctx, tt.disarmCommand, tt.targetSystem, tt.targetComponent)
				if disarmErr != nil {
					if !tt.wantErr {
						t.Errorf("ArmDisarm() failed: %v", disarmErr)
					}
					return
				}
			}
			s.Close()

		})
	}
}

func TestDroneAPI_DownloadMissions(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		targetSystem    uint8
		targetComponent uint8
		clientAddress   string
		outVersion      gomavlib.Version
		outSystemID     byte
		want            map[uint16]*ardupilotmega.MessageMissionItemInt
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
			got, gotErr := s.DownloadMissions(ctx, tt.targetSystem, tt.targetComponent)
			if gotErr == nil {
				if !tt.wantErr {
					t.Errorf("DownloadMissions() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("DownloadMissions() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if got != nil {
				t.Errorf("DownloadMissions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDroneAPI_StartMission(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		targetSystem    uint8
		targetComponent uint8
		wantErr         bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var s gomavlinkdroneapi.DroneAPI
			gotErr := s.StartMission(context.Background(), tt.targetSystem, tt.targetComponent)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("StartMission() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("StartMission() succeeded unexpectedly")
			}
		})
	}
}
