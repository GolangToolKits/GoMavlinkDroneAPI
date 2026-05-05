package gomavlinkdroneapi

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// ArmDisarm // command Arm = 1, command Disarm = 0, examples: targetSystem = 1, targetComponent: 1
func (s *DroneAPI) ArmDisarm(ctx context.Context, arm bool, targetSystem uint8, targetComponent uint8) error {
	var param1 float32 = 0 // Disarm
	action := "Disarm"
	if arm {
		param1 = 1 // Arm
		action = "Arm"
	} else {
		param1 = 0 // Disarm
		action = "Disarm"
	}

	armCmd := &common.MessageCommandLong{
		TargetSystem:    targetSystem,
		TargetComponent: targetComponent,
		Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
		Param1:          param1,
		// Param2: 21196, // ArduPilot force-arm magic number (use with caution!)
	}

	// 1. Send the command
	if err := s.drone.node.WriteMessageAll(armCmd); err != nil {
		return fmt.Errorf("failed to send %s command: %w", action, err)
	}
	log.Printf("%s command sent, waiting for acknowledgment...", action)

	// 2. Wait for ACK
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%s timed out: %w", action, ctx.Err())

		case evt, ok := <-s.drone.node.Events():
			if !ok {
				return errors.New("event channel closed")
			}

			frm, ok := evt.(*gomavlib.EventFrame)
			if !ok || frm.SystemID() != targetSystem || frm.ComponentID() != targetComponent {
				continue
			}

			if msg, ok := frm.Message().(*common.MessageCommandAck); ok {
				if msg.Command != common.MAV_CMD_COMPONENT_ARM_DISARM {
					continue
				}

				switch msg.Result {
				case common.MAV_RESULT_ACCEPTED:
					log.Printf("Success: Vehicle %sed!", action)
					return nil
				case common.MAV_RESULT_TEMPORARILY_REJECTED:
					return errors.New("arm rejected: pre-arm checks failed (check GPS/Battery)")
				case common.MAV_RESULT_DENIED:
					return errors.New("arm denied: check safety switch or flight mode")
				default:
					return fmt.Errorf("arm failed with result code: %v", msg.Result)
				}
			}
		}
	}
}

func (s *DroneAPI) SetMode(targetSystem uint8, targetComponent uint8, baseMode uint8, customMode uint32) error {
	cmd := &common.MessageCommandLong{
		TargetSystem: targetSystem,
		//
		TargetComponent: targetComponent, // Usually 1 for the autopilot
		Command:         common.MAV_CMD_DO_SET_MODE,
		Param1:          float32(baseMode),   //1
		Param2:          float32(customMode), //4
	}

	return s.drone.node.WriteMessageAll(cmd)
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
func (s *DroneAPI) Takeoff(ctx context.Context, altitude float32, targetSystem uint8, targetComponent uint8) error {
	takeOffCmd := &common.MessageCommandLong{
		TargetSystem:    targetSystem,
		TargetComponent: targetComponent,
		Command:         common.MAV_CMD_NAV_TAKEOFF,
		Param7:          altitude, // Param7 is typically Altitude for NAV_TAKEOFF
	}

	// 1. Send the command
	if err := s.drone.node.WriteMessageAll(takeOffCmd); err != nil {
		return fmt.Errorf("failed to send takeoff: %w", err)
	}
	log.Printf("Takeoff command (Alt: %.1fm) sent, waiting for ACK...", altitude)

	// 2. Wait for confirmation
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("takeoff timed out: %w", ctx.Err())

		case evt, ok := <-s.drone.node.Events():
			if !ok {
				return errors.New("event channel closed")
			}

			frm, ok := evt.(*gomavlib.EventFrame)
			if !ok || frm.SystemID() != targetSystem || frm.ComponentID() != targetComponent {
				continue
			}

			if msg, ok := frm.Message().(*common.MessageCommandAck); ok {
				if msg.Command != common.MAV_CMD_NAV_TAKEOFF {
					continue
				}

				switch msg.Result {
				case common.MAV_RESULT_ACCEPTED:
					log.Println("Takeoff accepted! Climbing...")
					return nil
				case common.MAV_RESULT_TEMPORARILY_REJECTED:
					return errors.New("takeoff rejected: is the drone armed and in GUIDED mode?")
				case common.MAV_RESULT_DENIED:
					return errors.New("takeoff denied: check safety limits or RC position")
				default:
					return fmt.Errorf("takeoff failed: result code %v", msg.Result)
				}
			}
		}
	}
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
func (s *DroneAPI) Land(ctx context.Context, targetSystem uint8, targetComponent uint8) error {
	landCmd := &common.MessageCommandLong{
		TargetSystem:    targetSystem,
		TargetComponent: targetComponent,
		Command:         common.MAV_CMD_NAV_LAND,
		// Param 1-7 are typically 0 for a basic land-at-current-position command
	}

	// 1. Send the command
	if err := s.drone.node.WriteMessageAll(landCmd); err != nil {
		return fmt.Errorf("failed to send land command: %w", err)
	}
	log.Println("Land command sent, waiting for acknowledgment...")

	// 2. Wait for ACK
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("landing request timed out: %w", ctx.Err())

		case evt, ok := <-s.drone.node.Events():
			if !ok {
				return errors.New("event channel closed")
			}

			frm, ok := evt.(*gomavlib.EventFrame)
			if !ok || frm.SystemID() != targetSystem || frm.ComponentID() != targetComponent {
				continue
			}

			if msg, ok := frm.Message().(*common.MessageCommandAck); ok {
				if msg.Command != common.MAV_CMD_NAV_LAND {
					continue
				}

				switch msg.Result {
				case common.MAV_RESULT_ACCEPTED:
					log.Println("Land command ACCEPTED. Vehicle is descending.")
					return nil
				case common.MAV_RESULT_TEMPORARILY_REJECTED:
					return errors.New("land rejected: check current flight mode")
				case common.MAV_RESULT_DENIED:
					return errors.New("land denied: command not allowed in current state")
				default:
					return fmt.Errorf("land failed: result code %v", msg.Result)
				}
			}
		}
	}
}

// AcknowledgeCommand commandToCheck:
// common.MAV_CMD_NAV_TAKEOFF,
// common.MAV_CMD_NAV_LAND,
// common.MAV_CMD_COMPONENT_ARM_DISARM
// func (s *DroneAPI) AcknowledgeCommand(commandToCheck common.MAV_CMD) bool {
// 	var rtn bool
// 	for evt := range s.drone.node.Events() {
// 		// 1. Check if the event is an incoming frame
// 		frm, ok := evt.(*gomavlib.EventFrame)
// 		if ok {
// 			// continue // Skip other events like channel connections/disconnections
// 			// 2. Check if the message is a Command Acknowledgment
// 			if ack, ok := frm.Message().(*common.MessageCommandAck); ok {
// 				if commandToCheck == ack.Command {
// 					rtn = true
// 					break
// 				}

// 				// 3. Match the command ID to verify which command this ACK belongs to
// 				// switch ack.Command {
// 				// case common.MAV_CMD_NAV_TAKEOFF:
// 				// 	log.Printf("Takeoff ACK Received! Result Code: %v\n", ack.Result)
// 				// case common.MAV_CMD_NAV_LAND:
// 				// 	log.Printf("Land ACK Received! Result Code: %v\n", ack.Result)
// 				// case common.MAV_CMD_COMPONENT_ARM_DISARM:
// 				// 	log.Printf("Arm/Disarm ACK Received! Result Code: %v\n", ack.Result)
// 				// default:
// 				// 	log.Printf("Received ACK for Command %d. Result: %v\n", ack.Command, ack.Result)
// 				// }

// 				// Interpret the result code
// 				// handleAckResult(ack.Result)
// 				// Helper function to interpret the MAV_RESULT enum
// 				// func handleAckResult(result common.MAV_RESULT) {
// 				// 	switch result {
// 				// 	case common.MAV_RESULT_ACCEPTED:
// 				// 		fmt.Println("🚀 Success: Command accepted and executed.")
// 				// 	case common.MAV_RESULT_TEMPORARILY_REJECTED:
// 				// 		fmt.Println("❌ Denied: Temporarily rejected (e.g., drone isn't armed yet or has no GPS lock).")
// 				// 	case common.MAV_RESULT_DENIED:
// 				// 		fmt.Println("❌ Denied: Command is invalid or refused by autopilot.")
// 				// 	case common.MAV_RESULT_UNSUPPORTED:
// 				// 		fmt.Println("❌ Denied: Autopilot doesn't support this command.")
// 				// 	default:
// 				// 		fmt.Printf("Notice: Other result code received (%d)\n", result)
// 				// 	}
// 				// }
// 			}
// 		}
// 	}
// 	return rtn
// }

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
func (s *DroneAPI) Move(targetSystem uint8, targetComponent uint8, moveMessage *common.MessageSetPositionTargetLocalNed) error {
	// 1. Explicitly set target IDs to ensure the message hits the right drone.
	moveMessage.TargetSystem = targetSystem
	moveMessage.TargetComponent = targetComponent

	// 2. Broadcast the message.
	// gomavlib expects the pointer (*common.MessageSetPositionTargetLocalNed).
	if err := s.drone.node.WriteMessageAll(moveMessage); err != nil {
		return fmt.Errorf("failed to send movement command: %w", err)
	}

	// 3. Log using the correct field names: Vx and Vy
	log.Printf("Sent Move Command: Vx=%.2f, Vy=%.2f to System %d",
		moveMessage.Vx, moveMessage.Vy, targetSystem)

	return nil
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
func (s *DroneAPI) ReturnHome(ctx context.Context, targetSystem uint8, targetComponent uint8) error {
	rtlCmd := &common.MessageCommandLong{
		TargetSystem:    targetSystem,
		TargetComponent: targetComponent,
		Command:         common.MAV_CMD_NAV_RETURN_TO_LAUNCH,
	}

	// 1. Send the command
	if err := s.drone.node.WriteMessageAll(rtlCmd); err != nil {
		return fmt.Errorf("failed to broadcast RTL: %w", err)
	}
	log.Println("RTL command sent, waiting for ACK...")

	// 2. Wait for the drone to acknowledge the mode change
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("RTL request timed out: %w", ctx.Err())

		case evt, ok := <-s.drone.node.Events():
			if !ok {
				return errors.New("event channel closed")
			}

			frm, ok := evt.(*gomavlib.EventFrame)
			if !ok || frm.SystemID() != targetSystem || frm.ComponentID() != targetComponent {
				continue
			}

			if msg, ok := frm.Message().(*common.MessageCommandAck); ok {
				if msg.Command != common.MAV_CMD_NAV_RETURN_TO_LAUNCH {
					continue
				}

				switch msg.Result {
				case common.MAV_RESULT_ACCEPTED:
					log.Println("RTL Accepted! Drone is coming home.")
					return nil
				case common.MAV_RESULT_TEMPORARILY_REJECTED:
					return errors.New("RTL rejected: check if home position is set")
				case common.MAV_RESULT_DENIED:
					return errors.New("RTL denied: drone refuses to switch modes")
				default:
					return fmt.Errorf("RTL failed with result: %v", msg.Result)
				}
			}
		}
	}
}
