package gomavlinkdroneapi_test

import (
	"context"
	"fmt"
	"testing"

	gomavlinkdroneapi "github.com/GolangToolKits/GoMavlinkDroneAPI"
	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

func TestDroneAPI_ArmDisarmTakeOffLand(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		//imputs for UDPClient node
		clientAddress string
		outVersion    gomavlib.Version
		outSystemID   byte
		// Named input parameters for target function.
		armCommand      bool
		disarmCommand   bool
		targetSystem    uint8
		targetComponent uint8

		checkConnect bool
		checkArm     bool
		wantErr      bool
		altitude     float32
		returnToHome bool
		moveMessage  *common.MessageSetPositionTargetLocalNed
	}{
		// TODO: Add test cases.
		{
			name:          "test 1",
			clientAddress: "127.0.0.1:5760",
			outVersion:    gomavlib.V2,
			outSystemID:   255,
			// commands
			checkConnect:    true,
			armCommand:      true,
			checkArm:        false,
			disarmCommand:   false,
			targetSystem:    1,
			targetComponent: 1,
			wantErr:         false,
			altitude:        1,
			returnToHome:    false,
			moveMessage: &common.MessageSetPositionTargetLocalNed{
				TargetSystem:    1, // Usually 1 for the first drone
				TargetComponent: 1, // Usually 1 for the main flight controller

				// MAV_FRAME_BODY_OFFSET_NED means directions are relative to the drone's nose
				CoordinateFrame: common.MAV_FRAME_BODY_OFFSET_NED,

				// 3527 = Ignore position & acceleration. ONLY look at velocity and yaw rate.
				TypeMask: 3527,

				Vx: 1.0, // Move Forward at 2.0 meters per second
				Vy: 0.0, // Do not strafe (Right)
				Vz: 0.0, // Do not change altitude (Down)
			},
		},
		{
			name:          "test 2",
			clientAddress: "127.0.0.1:5760",
			outVersion:    gomavlib.V2,
			outSystemID:   255,
			// commands
			checkConnect:    true,
			armCommand:      true,
			checkArm:        false,
			disarmCommand:   false,
			targetSystem:    1,
			targetComponent: 1,
			wantErr:         false,
			altitude:        1,
			returnToHome:    true,
			moveMessage: &common.MessageSetPositionTargetLocalNed{
				TargetSystem:    1, // Usually 1 for the first drone
				TargetComponent: 1, // Usually 1 for the main flight controller

				// MAV_FRAME_BODY_OFFSET_NED means directions are relative to the drone's nose
				CoordinateFrame: common.MAV_FRAME_BODY_OFFSET_NED,

				// 3527 = Ignore position & acceleration. ONLY look at velocity and yaw rate.
				TypeMask: 3527,

				Vx: 1.0, // Move Forward at 2.0 meters per second
				Vy: 0.0, // Do not strafe (Right)
				Vz: 0.0, // Do not change altitude (Down)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			// gotErr := s.ConnectTCPServer(tt.clientAddress, tt.outVersion, tt.outSystemID)
			// if gotErr != nil {
			// 	if !tt.wantErr {
			// 		t.Errorf("ConnectTCPClient() failed: %v", gotErr)
			// 	}
			// 	return
			// }
			if tt.checkConnect {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				con, err := s.IsDroneConnected(ctx)
				if !con {
					fmt.Print(err)
					t.Fatal("ConnectSerial() succeeded not connected")
				}
				//return
			}

			modeErr := s.SetMode(tt.targetSystem, tt.targetComponent, gomavlinkdroneapi.MODE_CUSTOM, gomavlinkdroneapi.MODE_GUIDED)
			if modeErr != nil {
				if !tt.wantErr {
					t.Errorf("SetMode() failed: %v", modeErr)
				}
				return
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			armErr := s.ArmDisarm(ctx, tt.armCommand, tt.targetSystem, tt.targetComponent)
			if armErr != nil {
				if !tt.wantErr {
					t.Errorf("ArmDisarm() failed: %v", armErr)
				}
				return
			}
			// if tt.checkArm {
			// 	armed := s.AcknowledgeCommand(common.MAV_CMD_COMPONENT_ARM_DISARM)
			// 	if !armed {
			// 		t.Fatal("Arm() failed")
			// 	}
			// 	return
			// }

			tkoffErr := s.Takeoff(ctx, tt.altitude, tt.targetSystem, tt.targetComponent)
			if tkoffErr != nil {
				if !tt.wantErr {
					t.Errorf("Takeoff() failed: %v", tkoffErr)
				}
				return
			}

			movErr := s.Move(tt.targetSystem, tt.targetComponent, tt.moveMessage)
			if movErr != nil {
				if !tt.wantErr {
					t.Errorf("Move() failed: %v", movErr)
				}
				return
			}

			if tt.returnToHome {
				homeErr := s.ReturnHome(ctx, tt.targetSystem, tt.targetComponent)
				if homeErr != nil {
					if !tt.wantErr {
						t.Errorf("ReturnHome() failed: %v", homeErr)
					}
					return
				}
			} else {
				lndErr := s.Land(ctx, tt.targetSystem, tt.targetComponent)
				if lndErr != nil {
					if !tt.wantErr {
						t.Errorf("Land() failed: %v", lndErr)
					}
					return
				}
			}

			disarmErr := s.ArmDisarm(ctx, tt.disarmCommand, tt.targetSystem, tt.targetComponent)
			if disarmErr != nil {
				if !tt.wantErr {
					t.Errorf("ArmDisarm() failed: %v", disarmErr)
				}
				return
			}
			s.Close()
			if tt.wantErr {
				t.Fatal("ArmDisarm() succeeded unexpectedly")
			}
		})
	}
}

func TestDroneAPI_Move(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		moveMessage     *common.MessageSetPositionTargetLocalNed
		wantErr         bool
		clientAddress   string
		outVersion      gomavlib.Version
		outSystemID     byte
		armCommand      bool
		disarmCommand   bool
		targetSystem    uint8
		targetComponent uint8
		takeOffCommand  *common.MessageCommandLong
		rtlCommand      *common.MessageCommandLong
		altitude        float32
	}{
		// TODO: Add test cases.

		{
			name:            "test 1",
			clientAddress:   "1.2.3.4:5600",
			outVersion:      gomavlib.V2,
			outSystemID:     10,
			armCommand:      true,
			disarmCommand:   true,
			targetSystem:    1,
			targetComponent: 1,
			altitude:        1,
			takeOffCommand: &common.MessageCommandLong{
				TargetSystem:    1,
				TargetComponent: 1,
				Command:         common.MAV_CMD_NAV_TAKEOFF,
				Param1:          0,   // Pitch (Minimum pitch for fixed wing, 0 for copters)
				Param4:          0,   // Yaw angle (0 = current heading)
				Param5:          0,   // Latitude (0 = current location)
				Param6:          0,   // Longitude (0 = current location)
				Param7:          1.0, // Altitude target in meters
			},
			moveMessage: &common.MessageSetPositionTargetLocalNed{
				TargetSystem:    1, // Usually 1 for the first drone
				TargetComponent: 1, // Usually 1 for the main flight controller

				// MAV_FRAME_BODY_OFFSET_NED means directions are relative to the drone's nose
				CoordinateFrame: common.MAV_FRAME_BODY_OFFSET_NED,

				// 3527 = Ignore position & acceleration. ONLY look at velocity and yaw rate.
				TypeMask: 3527,

				Vx: 1.0, // Move Forward at 2.0 meters per second
				Vy: 0.0, // Do not strafe (Right)
				Vz: 0.0, // Do not change altitude (Down)
			},
			rtlCommand: &common.MessageCommandLong{
				TargetSystem:    1, // System ID of your drone
				TargetComponent: 1, // Component ID of your drone's flight controller
				Command:         common.MAV_CMD_NAV_RETURN_TO_LAUNCH,
				Confirmation:    1, // 0 for first transmission, 1 for confirmation
				Param1:          0, // Unused for RTL
				Param2:          0, // Unused for RTL
				Param3:          0, // Unused for RTL
				Param4:          0, // Unused for RTL
				Param5:          0, // Unused for RTL
				Param6:          0, // Unused for RTL
				Param7:          0, // Unused for RTL
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
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			armErr := s.ArmDisarm(ctx, tt.armCommand, tt.targetSystem, tt.targetComponent)
			if armErr != nil {
				if !tt.wantErr {
					t.Errorf("ArmDisarm() failed: %v", armErr)
				}
				return
			}
			tkoffErr := s.Takeoff(ctx, tt.altitude, tt.targetSystem, tt.targetComponent)
			if tkoffErr != nil {
				if !tt.wantErr {
					t.Errorf("Takeoff() failed: %v", tkoffErr)
				}
				return
			}

			gotErr := s.Move(tt.targetSystem, tt.targetComponent, tt.moveMessage)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Move() failed: %v", gotErr)
				}
				return
			}

			homeErr := s.ReturnHome(ctx, tt.targetSystem, tt.targetComponent)
			if homeErr != nil {
				if !tt.wantErr {
					t.Errorf("ReturnHome() failed: %v", homeErr)
				}
				return
			}
			disarmErr := s.ArmDisarm(ctx, tt.disarmCommand, tt.targetSystem, tt.targetComponent)
			if disarmErr != nil {
				if !tt.wantErr {
					t.Errorf("ArmDisarm() failed: %v", disarmErr)
				}
				return
			}
			s.Close()
			if tt.wantErr {
				t.Fatal("Move() succeeded unexpectedly")
			}
		})
	}
}

func TestDroneAPI_Takeoff(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		takeOffCommand  *common.MessageCommandLong
		wantErr         bool
		altitude        float32
		targetSystem    uint8
		targetComponent uint8
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var s gomavlinkdroneapi.DroneAPI
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			gotErr := s.Takeoff(ctx, tt.altitude, tt.targetSystem, tt.targetComponent)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Takeoff() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Takeoff() succeeded unexpectedly")
			}
		})
	}
}

func TestDroneAPI_Land(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		landCommand     *common.MessageCommandLong
		wantErr         bool
		clientAddress   string
		outVersion      gomavlib.Version
		outSystemID     byte
		armCommand      bool
		disarmCommand   bool
		targetSystem    uint8
		targetComponent uint8
		takeOffCommand  *common.MessageCommandLong
		altitude        float32
	}{
		// TODO: Add test cases.
		{
			name:            "test 1",
			clientAddress:   "1.2.3.4:5600",
			outVersion:      gomavlib.V2,
			outSystemID:     10,
			armCommand:      true,
			disarmCommand:   true,
			targetSystem:    1,
			targetComponent: 1,
			altitude:        1,
			takeOffCommand: &common.MessageCommandLong{
				TargetSystem:    1,
				TargetComponent: 1,
				Command:         common.MAV_CMD_NAV_TAKEOFF,
				Param1:          0,   // Pitch (Minimum pitch for fixed wing, 0 for copters)
				Param4:          0,   // Yaw angle (0 = current heading)
				Param5:          0,   // Latitude (0 = current location)
				Param6:          0,   // Longitude (0 = current location)
				Param7:          1.0, // Altitude target in meters
			},
			landCommand: &common.MessageCommandLong{
				TargetSystem:    1,                       // System ID of your drone
				TargetComponent: 1,                       // Component ID of your drone
				Command:         common.MAV_CMD_NAV_LAND, // Command ID 21
				Param1:          0,                       // Abort Alt (0 = use autopilot default behavior)
				Param2:          0,                       // Precision land mode (0 = normal land)
				Param4:          0,                       // Yaw angle (0 = ignore/current heading)
				Param5:          0,                       // Latitude (0 = current location)
				Param6:          0,                       // Longitude (0 = current location)
				Param7:          0,                       // Altitude (0 = ground level)
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
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			armErr := s.ArmDisarm(ctx, tt.armCommand, tt.targetSystem, tt.targetComponent)
			if armErr != nil {
				if !tt.wantErr {
					t.Errorf("ArmDisarm() failed: %v", armErr)
				}
				return
			}
			tkoffErr := s.Takeoff(ctx, tt.altitude, tt.targetSystem, tt.targetComponent)
			if tkoffErr != nil {
				if !tt.wantErr {
					t.Errorf("Takeoff() failed: %v", tkoffErr)
				}
				return
			}
			gotErr := s.Land(ctx, tt.targetSystem, tt.targetComponent)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Land() failed: %v", gotErr)
				}
				return
			}
			disarmErr := s.ArmDisarm(ctx, tt.disarmCommand, tt.targetSystem, tt.targetComponent)
			if disarmErr != nil {
				if !tt.wantErr {
					t.Errorf("ArmDisarm() failed: %v", disarmErr)
				}
				return
			}
			s.Close()
			if tt.wantErr {
				t.Fatal("Land() succeeded unexpectedly")
			}
		})
	}
}

func TestDroneAPI_ReturnHome(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		rtlCommand      *common.MessageCommandLong
		wantErr         bool
		targetSystem    uint8
		targetComponent uint8
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var s gomavlinkdroneapi.DroneAPI
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			gotErr := s.ReturnHome(ctx, tt.targetSystem, tt.targetComponent)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ReturnHome() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ReturnHome() succeeded unexpectedly")
			}
		})
	}
}
