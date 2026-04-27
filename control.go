package gomavlinkdroneapi

import (
	"log"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// ArmDisarm // command Arm = 1, command Disarm = 0, examples: targetSystem = 1, targetComponent: 1
func (s *DroneAPI) ArmDisarm(command float32, targetSystem uint8, targetComponent uint8) error {
	var rtn error
	armCmd := &common.MessageCommandLong{
		//TargetSystem:    1, // System ID of your drone
		TargetSystem: targetSystem, // System ID of your drone
		//TargetComponent: 1, // Component ID of your drone
		TargetComponent: targetComponent, // Component ID of your drone
		Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
		Param1:          command, // 1 = Arm, 0 = Disarm
	}
	rtn = s.drone.node.WriteMessageAll(armCmd)
	switch command {
	case 1:
		log.Println("Arm command sent.")
	case 0:
		log.Println("Disarm command sent.")
	}

	// Wait a moment for motors to fully spin up
	time.Sleep(3 * time.Second)
	return rtn
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
	rtn := s.drone.node.WriteMessageAll(takeOffCommand)
	log.Println("Takeoff command sent for 1 meters altitude.")
	return rtn
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
	rtn := s.drone.node.WriteMessageAll(landCommand)
	log.Println("Land command sent for 0 meters ground.")
	return rtn
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
