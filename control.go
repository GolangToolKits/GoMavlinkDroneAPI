package gomavlinkdroneapi

import (
	"log"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// ArmDisarm // command Arm = 1, command Disarm = 0, examples: targetSystem = 1, targetComponent: 1
func (s *DroneAPI) ArmDisarm(command float32, targetSystem uint8, targetComponent uint8) error {
	//var rtn error
	armCmd := &common.MessageCommandLong{
		//TargetSystem:    1, // System ID of your drone
		TargetSystem: targetSystem, // System ID of your drone
		//TargetComponent: 1, // Component ID of your drone
		TargetComponent: targetComponent, // Component ID of your drone
		Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
		Param1:          command, // 1 = Arm, 0 = Disarm
	}
	err := s.drone.node.WriteMessageAll(armCmd)
	switch command {
	case 1:
		log.Println("Arm command sent.")
	case 0:
		log.Println("Disarm command sent.")
	}

	// Wait a moment for motors to fully spin up
	//time.Sleep(3 * time.Second)
	return err
}

// Takeoff sends takeoff
// example
//
//	&common.MessageCommandLong{
//			TargetSystem:    1,
//			TargetComponent: 1,
//			Command:         common.MAV_CMD_NAV_TAKEOFF,
//			Param1:         0,   // Pitch (Minimum pitch for fixed wing, 0 for copters)
//			Param4:         0,   // Yaw angle (0 = current heading)
//			Param5:         0,   // Latitude (0 = current location)
//			Param6:         0,   // Longitude (0 = current location)
//			Param7:         1.0, // Altitude target in meters
//		}
func (s *DroneAPI) Takeoff(takeOffCommand *common.MessageCommandLong) error {
	err := s.drone.node.WriteMessageAll(takeOffCommand)
	log.Println("Takeoff command sent for 1 meters altitude.")
	return err
}

// Land sends land command
// example
//
//	landCmd := &common.MessageCommandLong{
//			TargetSystem:    1, // System ID of your drone
//			TargetComponent: 1, // Component ID of your drone
//			Command:         common.MAV_CMD_NAV_LAND, // Command ID 21
//			Param1:         0, // Abort Alt (0 = use autopilot default behavior)
//			Param2:         0, // Precision land mode (0 = normal land)
//			Param4:         0, // Yaw angle (0 = ignore/current heading)
//			Param5:         0, // Latitude (0 = current location)
//			Param6:         0, // Longitude (0 = current location)
//			Param7:         0, // Altitude (0 = ground level)
//		}
func (s *DroneAPI) Land(landCommand *common.MessageCommandLong) error {
	err := s.drone.node.WriteMessageAll(landCommand)
	log.Println("Land command sent for 0 meters ground.")
	return err
}

// AcknowledgeCommand commandToCheck:
// common.MAV_CMD_NAV_TAKEOFF,
// common.MAV_CMD_NAV_LAND,
// common.MAV_CMD_COMPONENT_ARM_DISARM
func (s *DroneAPI) AcknowledgeCommand(commandToCheck common.MAV_CMD) bool {
	var rtn bool
	for evt := range s.drone.node.Events() {
		// 1. Check if the event is an incoming frame
		frm, ok := evt.(*gomavlib.EventFrame)
		if ok {
			// continue // Skip other events like channel connections/disconnections
			// 2. Check if the message is a Command Acknowledgment
			if ack, ok := frm.Message().(*common.MessageCommandAck); ok {
				if commandToCheck == ack.Command {
					rtn = true
					break
				}

				// 3. Match the command ID to verify which command this ACK belongs to
				// switch ack.Command {
				// case common.MAV_CMD_NAV_TAKEOFF:
				// 	log.Printf("Takeoff ACK Received! Result Code: %v\n", ack.Result)
				// case common.MAV_CMD_NAV_LAND:
				// 	log.Printf("Land ACK Received! Result Code: %v\n", ack.Result)
				// case common.MAV_CMD_COMPONENT_ARM_DISARM:
				// 	log.Printf("Arm/Disarm ACK Received! Result Code: %v\n", ack.Result)
				// default:
				// 	log.Printf("Received ACK for Command %d. Result: %v\n", ack.Command, ack.Result)
				// }

				// Interpret the result code
				// handleAckResult(ack.Result)
				// Helper function to interpret the MAV_RESULT enum
				// func handleAckResult(result common.MAV_RESULT) {
				// 	switch result {
				// 	case common.MAV_RESULT_ACCEPTED:
				// 		fmt.Println("🚀 Success: Command accepted and executed.")
				// 	case common.MAV_RESULT_TEMPORARILY_REJECTED:
				// 		fmt.Println("❌ Denied: Temporarily rejected (e.g., drone isn't armed yet or has no GPS lock).")
				// 	case common.MAV_RESULT_DENIED:
				// 		fmt.Println("❌ Denied: Command is invalid or refused by autopilot.")
				// 	case common.MAV_RESULT_UNSUPPORTED:
				// 		fmt.Println("❌ Denied: Autopilot doesn't support this command.")
				// 	default:
				// 		fmt.Printf("Notice: Other result code received (%d)\n", result)
				// 	}
				// }
			}
		}
	}
	return rtn
}

//Move
//Examples
// moveMessage := &common.MessageSetPositionTargetLocalNed{
// 		TargetSystem:    1, // Usually 1 for the first drone
// 		TargetComponent: 1, // Usually 1 for the main flight controller

// 		// MAV_FRAME_BODY_OFFSET_NED means directions are relative to the drone's nose
// 		CoordinateFrame: common.MAV_FRAME_BODY_OFFSET_NED,

// 		// 3527 = Ignore position & acceleration. ONLY look at velocity and yaw rate.
// 		TypeMask: 3527,

//		Vx: 2.0, // Move Forward at 2.0 meters per second
//		Vy: 0.0, // Do not strafe (Right)
//		Vz: 0.0, // Do not change altitude (Down)
//	}
//
// !!!!!!!!!!!!  Notice  !!!!!!!!!!!
// Actively broadcasting movement commands at a rate of 5Hz
// ticker := time.NewTicker(200 * time.Millisecond)
// This most be done by the calling app:
// A WebSocket app
func (s *DroneAPI) Move(moveMessage *common.MessageSetPositionTargetLocalNed) error {
	// Arming & Mode: This code strictly streams the request to move. For the drone to actually physicalize
	// this request, you must independently command the vehicle to Arm its motors and change its flight mode
	// to Guided (Ardupilot) or Offboard (PX4).
	// You can do this via a manual transmitter or by writing automated commands
	// through node.WriteMessageAll as well.
	err := s.drone.node.WriteMessageAll(moveMessage)
	if err != nil {
		log.Println("Error broadcasting movement command:", err)
	} else {
		log.Println("Sent: Moving Forward at 2m/s")
	}
	return err
}

// Returnhome send vehicle back to launch position
// Example
//
//	rtlCommand := &common.MessageCommandLong{
//			TargetSystem:    1, // System ID of your drone
//			TargetComponent: 1, // Component ID of your drone's flight controller
//			Command:         common.MAV_CMD_NAV_RETURN_TO_LAUNCH,
//			Confirmation:    1, // 0 for first transmission, 1 for confirmation
//			Param1:          0, // Unused for RTL
//			Param2:          0, // Unused for RTL
//			Param3:          0, // Unused for RTL
//			Param4:          0, // Unused for RTL
//			Param5:          0, // Unused for RTL
//			Param6:          0, // Unused for RTL
//			Param7:          0, // Unused for RTL
//		}
func (s *DroneAPI) ReturnHome(rtlCommand *common.MessageCommandLong) error {
	err := s.drone.node.WriteMessageAll(rtlCommand)
	if err != nil {
		log.Println("Error sending RTL command:", err)
	} else {
		log.Println("Return To Launch command sent successfully!")
	}
	return err
}
