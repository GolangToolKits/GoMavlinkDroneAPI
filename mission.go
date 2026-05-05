package gomavlinkdroneapi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// UploadMission
// examples
//
//	var missionItems = []common.MessageMissionItemInt{
//			{
//				TargetSystem:    1,
//				TargetComponent: 1,
//				Seq:             0,
//				Frame:           common.MAV_FRAME_GLOBAL_RELATIVE_ALT,
//				Command:         common.MAV_CMD_NAV_TAKEOFF,
//				Current:         1,
//				Autocontinue:    1,
//				Param7:          10, // Take off to 10 meters
//			},
//			{
//				TargetSystem:    1,
//				TargetComponent: 1,
//				Seq:             1,
//				Frame:           common.MAV_FRAME_GLOBAL_RELATIVE_ALT,
//				Command:         common.MAV_CMD_NAV_WAYPOINT,
//				Current:         0,
//				Autocontinue:    1,
//				X:               473977418, // Latitude * 1e7
//				Y:               85443833,  // Longitude * 1e7
//				Param7:          15,        // Altitude (meters)
//			},
//		}
const uploadTimeout = 10 * time.Second

func (s *DroneAPI) UploadMission(ctx context.Context, missionItems []common.MessageMissionItemInt, targetSystem uint8, targetComponent uint8) error {
	totalItems := uint16(len(missionItems))
	if totalItems == 0 {
		return errors.New("cannot upload an empty mission")
	}

	// 1. Send Count
	s.drone.node.WriteMessageAll(&common.MessageMissionCount{
		TargetSystem:    targetSystem,
		TargetComponent: targetComponent,
		Count:           totalItems,
		MissionType:     common.MAV_MISSION_TYPE_MISSION,
	})

	// Use a deadline for the entire operation or per-packet
	expiry := time.Now().Add(uploadTimeout)

	for {
		if time.Now().After(expiry) {
			return errors.New("mission upload timed out")
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case evt, ok := <-s.drone.node.Events():
			if !ok {
				return errors.New("node events channel closed")
			}

			frm, ok := evt.(*gomavlib.EventFrame)
			if !ok || frm.SystemID() != targetSystem {
				continue
			}

			// Handle different request types (Standard and Int)
			var requestedSeq uint16
			isRequest := false

			switch msg := frm.Message().(type) {
			case *common.MessageMissionRequestInt:
				requestedSeq = msg.Seq
				isRequest = true
			case *common.MessageMissionRequest:
				requestedSeq = msg.Seq
				isRequest = true
			case *common.MessageMissionAck:
				if msg.Type == common.MAV_MISSION_ACCEPTED {
					log.Println("Mission accepted!")
					return nil
				}
				return fmt.Errorf("mission rejected: %v", msg.Type)
			}

			if isRequest {
				if int(requestedSeq) >= len(missionItems) {
					return fmt.Errorf("drone requested out-of-bounds seq %d", requestedSeq)
				}

				// Prepare item and ensure IDs match
				item := missionItems[requestedSeq]
				item.TargetSystem = targetSystem
				item.TargetComponent = targetComponent
				item.Seq = requestedSeq // Ensure sequence matches request

				s.drone.node.WriteMessageTo(frm.Channel, &item)
				expiry = time.Now().Add(uploadTimeout) // Refresh timeout
			}
		}
	}
}

// Start Drone mission  !!important!! Drone must be armed before calling this method.
func (s *DroneAPI) StartMission(ctx context.Context, targetSystem uint8, targetComponent uint8) error {
	// 1. Send the command to switch to AUTO mode
	err := s.drone.node.WriteMessageAll(&common.MessageCommandLong{
		TargetSystem:    targetSystem,
		TargetComponent: targetComponent,
		Command:         common.MAV_CMD_DO_SET_MODE,
		Param1:          float32(common.MAV_MODE_FLAG_CUSTOM_MODE_ENABLED),
		Param2:          3, // ArduCopter AUTO; use 4 for PX4
	})
	if err != nil {
		return fmt.Errorf("failed to send mode command: %w", err)
	}

	// 2. Wait for COMMAND_ACK
	for {
		select {
		case <-ctx.Done():
			return errors.New("timeout waiting for command acknowledgment")
		case evt := <-s.drone.node.Events():
			if frm, ok := evt.(*gomavlib.EventFrame); ok {
				if msg, ok := frm.Message().(*common.MessageCommandAck); ok {
					// Check if this ACK is for our specific command
					if msg.Command == common.MAV_CMD_DO_SET_MODE {
						if msg.Result == common.MAV_RESULT_ACCEPTED {
							log.Println("Mission started successfully!")
							return nil
						}
						return fmt.Errorf("mission start rejected: %v", msg.Result)
					}
				}
			}
		}
	}
}

// MonitorMission listens for mission progress and sends RTL when the final waypoint is reached
func (s *DroneAPI) MonitorMission(ctx context.Context, totalItems uint16, targetSystem uint8) (bool, error) {
	lastSeq := totalItems - 1
	var currentTarget uint16
	s.drone.node.WriteMessageAll(&common.MessageCommandLong{
		TargetSystem:    1,
		TargetComponent: 1,
		Command:         common.MAV_CMD_SET_MESSAGE_INTERVAL,
		Param1:          46,     //42,     // Message ID for MISSION_CURRENT
		Param2:          500000, // Interval in microseconds
	})

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()

		case evt, ok := <-s.drone.node.Events():
			if !ok {
				return false, fmt.Errorf("node events channel closed")
			}

			frm, ok := evt.(*gomavlib.EventFrame)
			if !ok || frm.SystemID() != targetSystem {
				continue
			}

			log.Printf("Received message of type %T", frm.Message())

			switch msg := frm.Message().(type) {

			case *common.MessageMissionCurrent:
				// Track which item the drone is currently flying TOWARD
				if msg.Seq != currentTarget {
					currentTarget = msg.Seq
					fmt.Printf("Drone is now navigating to waypoint #%d\n", currentTarget)
				}

			case *common.MessageMissionItemReached:
				// Track which item the drone has ARRIVED at
				fmt.Printf("Successfully reached waypoint #%d/%d\n", msg.Seq, lastSeq)

				if msg.Seq == MISSION_COMPLETED {
					fmt.Println("Final waypoint reached. Sending RTL...")
					err := s.drone.node.WriteMessageAll(&common.MessageCommandLong{
						TargetSystem: targetSystem,
						Command:      common.MAV_CMD_NAV_RETURN_TO_LAUNCH,
					})
					return true, err
				}
			}
		}
	}
}

const (
	packetTimeout = 2 * time.Second
	maxRetries    = 5
)

func (s *DroneAPI) DownloadMissions(ctx context.Context, targetSystem uint8, targetComponent uint8) ([]*common.MessageMissionItemInt, error) {
	log.Println("Requesting mission list...")

	var items []*common.MessageMissionItemInt
	var total uint16
	var commsChannel *gomavlib.Channel // Track the specific channel the drone is on

	current := uint16(0)
	countReceived := false
	retries := 0

	// 1. Initial Request
	s.drone.node.WriteMessageAll(&common.MessageMissionRequestList{
		TargetSystem:    targetSystem,
		TargetComponent: targetComponent,
		MissionType:     common.MAV_MISSION_TYPE_MISSION,
	})

	timer := time.NewTimer(packetTimeout)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		case <-timer.C:
			if retries >= maxRetries {
				return nil, fmt.Errorf("download failed: max retries (%d) reached at item %d", maxRetries, current)
			}
			retries++
			log.Printf("Retry %d/%d for item #%d...", retries, maxRetries, current)

			if !countReceived {
				s.drone.node.WriteMessageAll(&common.MessageMissionRequestList{
					TargetSystem: targetSystem, TargetComponent: targetComponent,
					MissionType: common.MAV_MISSION_TYPE_MISSION,
				})
			} else if commsChannel != nil {
				s.drone.node.WriteMessageTo(commsChannel, &common.MessageMissionRequestInt{
					TargetSystem: targetSystem, TargetComponent: targetComponent,
					Seq: current, MissionType: common.MAV_MISSION_TYPE_MISSION,
				})
			}
			timer.Reset(packetTimeout)

		case evt, ok := <-s.drone.node.Events():
			if !ok {
				return nil, errors.New("node events channel closed")
			}
			frm, ok := evt.(*gomavlib.EventFrame)
			if !ok || frm.SystemID() != targetSystem {
				continue
			}

			switch msg := frm.Message().(type) {
			case *common.MessageMissionCount:
				if msg.MissionType != common.MAV_MISSION_TYPE_MISSION || countReceived {
					continue
				}

				commsChannel = frm.Channel
				total = msg.Count
				countReceived = true
				items = make([]*common.MessageMissionItemInt, total) // Fixed type

				if total == 0 {
					s.drone.node.WriteMessageTo(frm.Channel, &common.MessageMissionAck{
						TargetSystem: targetSystem, TargetComponent: targetComponent,
						Type: common.MAV_MISSION_ACCEPTED, MissionType: common.MAV_MISSION_TYPE_MISSION,
					})
					return []*common.MessageMissionItemInt{}, nil
				}

				retries = 0
				s.drone.node.WriteMessageTo(frm.Channel, &common.MessageMissionRequestInt{
					TargetSystem: targetSystem, TargetComponent: targetComponent,
					Seq: 0, MissionType: common.MAV_MISSION_TYPE_MISSION,
				})
				resetTimer(timer, packetTimeout)

			case *common.MessageMissionItemInt:
				if msg.MissionType != common.MAV_MISSION_TYPE_MISSION || msg.Seq != current {
					continue
				}

				items[msg.Seq] = msg
				log.Printf("Received item #%d/%d", msg.Seq+1, total)

				current++
				retries = 0

				if current == total {
					s.drone.node.WriteMessageTo(frm.Channel, &common.MessageMissionAck{
						TargetSystem: targetSystem, TargetComponent: targetComponent,
						Type: common.MAV_MISSION_ACCEPTED, MissionType: common.MAV_MISSION_TYPE_MISSION,
					})
					return items, nil
				}

				s.drone.node.WriteMessageTo(frm.Channel, &common.MessageMissionRequestInt{
					TargetSystem: targetSystem, TargetComponent: targetComponent,
					Seq: current, MissionType: common.MAV_MISSION_TYPE_MISSION,
				})
				resetTimer(timer, packetTimeout)
			}
		}
	}
}

// Helper to safely reset timers
func resetTimer(t *time.Timer, d time.Duration) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
	t.Reset(d)
}

// Prepare the MAV_CMD_DO_SET_MODE command
// Note: Mode numbers depend heavily on whether you use PX4 or ArduPilot
// PX4 HOLD mode is typically custom_mode = 4 (or send MAV_CMD_NAV_LOITER_UNLIM)
//
//	command := &common.MessageCommandLong{
//		TargetSystem:    1, // System ID of the drone
//		TargetComponent: 1, // Component ID of the drone
//		Command:         common.MAV_CMD_DO_REPOSITION,
//		Param1:          float32(common.MAV_MODE_FLAG_CUSTOM_MODE_ENABLED),
//		Param2:          4, // Example custom mode (HOLD in PX4)
//	}
//
//	func (s *DroneAPI) OverrideMissionAndHover(command *common.MessageCommandLong) error {
//		err := s.drone.node.WriteMessageAll(command)
//		log.Println("Sent override command: Hovering initiated.")
//		return err
//	}
func (s *DroneAPI) OverrideMissionAndHover(ctx context.Context, command *common.MessageCommandLong) error {
	var correctCommand bool
	if command.Command == common.MAV_CMD_DO_REPOSITION {
		correctCommand = true
	} else if !correctCommand && command.Command != common.MAV_CMD_NAV_LOITER_UNLIM {
		return fmt.Errorf("invalid command ID %d", command.Command)
	}

	// Use a ticker for retries. MAVLink commands often need 2-3 tries on noisy links.
	retryTicker := time.NewTicker(10 * time.Second)
	defer retryTicker.Stop()

	// Initial send
	s.drone.node.WriteMessageAll(command)
	log.Printf("Sending override command %d...", command.Command)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-retryTicker.C:
			log.Printf("Retrying command %d...", command.Command)
			s.drone.node.WriteMessageAll(command)

		case evt, ok := <-s.drone.node.Events():
			if !ok {
				return errors.New("node closed")
			}

			frm, ok := evt.(*gomavlib.EventFrame)
			// Check SystemID, but be careful with ComponentID.
			// A drone (CompID 1) often sends ACKs from its Autopilot component.
			if !ok || frm.SystemID() != command.TargetSystem {
				continue
			}

			if msg, ok := frm.Message().(*common.MessageCommandAck); ok {
				if msg.Command != command.Command {
					continue
				}

				switch msg.Result {
				case common.MAV_RESULT_ACCEPTED:
					log.Println("Command ACCEPTED")
					return nil
				case common.MAV_RESULT_IN_PROGRESS:
					log.Println("Command in progress...")
					continue // Keep waiting, don't retry yet
				case common.MAV_RESULT_TEMPORARILY_REJECTED:
					return errors.New("temporarily rejected: check flight mode/arming")
				case common.MAV_RESULT_DENIED:
					return errors.New("command denied")
				default:
					return fmt.Errorf("command failed: result %v", msg.Result)
				}
			}
		}
	}
}
