package gomavlinkdroneapi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
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
const uploadTimeout = 3 * time.Second

func (s *DroneAPI) UploadMission(ctx context.Context, missionItems []ardupilotmega.MessageMissionItemInt, targetSystem uint8, targetComponent uint8) error {
	totalItems := uint16(len(missionItems))
	if totalItems == 0 {
		return errors.New("cannot upload an empty mission")
	}

	log.Printf("Starting upload of %d waypoints...\n", totalItems)

	// 1. Initiate upload by sending MISSION_COUNT
	err := s.drone.node.WriteMessageAll(&common.MessageMissionCount{
		TargetSystem:    targetSystem,
		TargetComponent: targetComponent,
		Count:           totalItems,
		MissionType:     common.MAV_MISSION_TYPE_MISSION,
	})
	if err != nil {
		return fmt.Errorf("failed to send mission count: %w", err)
	}

	timer := time.NewTimer(uploadTimeout)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Mission upload aborted by external signal.")
			return ctx.Err()

		case <-timer.C:
			log.Println("Timeout waiting for drone to request items. Aborting upload.")
			return errors.New("mission upload timed out")

		case evt, ok := <-s.drone.node.Events():
			if !ok {
				return errors.New("node events channel closed unexpectedly")
			}

			frm, ok := evt.(*gomavlib.EventFrame)
			if !ok {
				continue
			}

			// Filter: Only process messages coming from our targeted drone
			if frm.SystemID() != targetSystem || frm.ComponentID() != targetComponent {
				continue
			}

			switch msg := frm.Message().(type) {

			// Drone is asking for a specific item using integer coordinates
			case *common.MessageMissionRequestInt:
				if msg.Seq >= totalItems {
					log.Printf("Drone requested out-of-bounds sequence: %d", msg.Seq)
					continue
				}

				itemToSend := missionItems[msg.Seq]
				itemToSend.TargetSystem = targetSystem
				itemToSend.TargetComponent = targetComponent

				s.drone.node.WriteMessageTo(frm.Channel, &itemToSend)
				log.Printf("Uploaded waypoint #%d", msg.Seq)

				// Reset the timer since the drone is actively communicating
				timer.Reset(uploadTimeout)

			// Handshake complete
			case *common.MessageMissionAck:
				if msg.Type == common.MAV_MISSION_ACCEPTED {
					log.Println("Mission uploaded and ACCEPTED by the drone!")
					return nil
				}

				log.Printf("Mission rejected by drone! Code: %d", msg.Type)
				return fmt.Errorf("drone rejected mission with code %d", msg.Type)
			}
		}
	}
}

const requestTimeout = 2 * time.Second

func (s *DroneAPI) DownloadMissions(ctx context.Context, targetSystem uint8, targetComponent uint8) ([]*ardupilotmega.MessageMissionItemInt, error) {
	log.Println("Requesting mission list from drone...")

	err := s.drone.node.WriteMessageAll(&ardupilotmega.MessageMissionRequestList{
		TargetSystem:    targetSystem,
		TargetComponent: targetComponent,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send mission request list: %w", err)
	}

	var missionItems []*ardupilotmega.MessageMissionItemInt
	var totalItems uint16
	currentRequested := uint16(0)
	countReceived := false

	// Timer to handle dropped MAVLink packets
	timer := time.NewTimer(requestTimeout)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Mission download aborted by external signal.")
			return nil, ctx.Err()

		case <-timer.C:
			// Packet loss occurred. Request the current index again.
			log.Printf("Timeout waiting for item #%d. Retrying request...", currentRequested)
			if !countReceived {
				s.drone.node.WriteMessageAll(&ardupilotmega.MessageMissionRequestList{
					TargetSystem:    targetSystem,
					TargetComponent: targetComponent,
				})
			} else {
				s.drone.node.WriteMessageAll(&ardupilotmega.MessageMissionRequestInt{
					TargetSystem:    targetSystem,
					TargetComponent: targetComponent,
					Seq:             currentRequested,
				})
			}
			timer.Reset(requestTimeout)

		case evt, ok := <-s.drone.node.Events():
			if !ok {
				return nil, errors.New("node events channel closed unexpectedly")
			}

			frm, ok := evt.(*gomavlib.EventFrame)
			if !ok {
				continue
			}

			// 1. Strict filtering: Ignore messages not sent by our targeted drone
			if frm.SystemID() != targetSystem || frm.ComponentID() != targetComponent {
				continue
			}

			switch msg := frm.Message().(type) {

			case *ardupilotmega.MessageMissionCount:
				if countReceived {
					continue // Ignore duplicate count signals
				}

				totalItems = msg.Count
				log.Printf("Drone reported %d items. Starting download...", totalItems)

				if totalItems == 0 {
					return nil, nil
				}

				countReceived = true
				missionItems = make([]*ardupilotmega.MessageMissionItemInt, totalItems)

				// Request first item
				currentRequested = 0
				s.drone.node.WriteMessageTo(frm.Channel, &ardupilotmega.MessageMissionRequestInt{
					TargetSystem:    targetSystem,
					TargetComponent: targetComponent,
					Seq:             currentRequested,
				})
				timer.Reset(requestTimeout)

			case *ardupilotmega.MessageMissionItemInt:
				if !countReceived || msg.Seq >= totalItems {
					continue // Out of bounds or premature item
				}

				// Only process if it is the specific item we asked for
				if msg.Seq == currentRequested && missionItems[msg.Seq] == nil {
					log.Printf("Received item #%d", msg.Seq)
					missionItems[msg.Seq] = msg
					currentRequested++

					// Check if download is complete
					if currentRequested == totalItems {
						log.Println("All items received successfully!")

						s.drone.node.WriteMessageTo(frm.Channel, &ardupilotmega.MessageMissionAck{
							TargetSystem:    targetSystem,
							TargetComponent: targetComponent,
							Type:            ardupilotmega.MAV_MISSION_ACCEPTED,
						})
						return missionItems, nil
					}

					// Request the next item
					s.drone.node.WriteMessageTo(frm.Channel, &ardupilotmega.MessageMissionRequestInt{
						TargetSystem:    targetSystem,
						TargetComponent: targetComponent,
						Seq:             currentRequested,
					})
					timer.Reset(requestTimeout)
				}
			}
		}
	}
}

// Prepare the MAV_CMD_DO_SET_MODE command
// Note: Mode numbers depend heavily on whether you use PX4 or ArduPilot
// PX4 HOLD mode is typically custom_mode = 4 (or send MAV_CMD_NAV_LOITER_UNLIM)
//
//	command := &common.MessageCommandLong{
//		TargetSystem:    1, // System ID of the drone
//		TargetComponent: 1, // Component ID of the drone
//		Command:         common.MAV_CMD_DO_SET_MODE,
//		Param1:          float32(common.MAV_MODE_FLAG_CUSTOM_MODE_ENABLED),
//		Param2:          4, // Example custom mode (HOLD in PX4)
//	}
//
//	func (s *DroneAPI) OverrideMissionAndHover(command *common.MessageCommandLong) error {
//		err := s.drone.node.WriteMessageAll(command)
//		log.Println("Sent override command: Hovering initiated.")
//		return err
//	}
const commandTimeout = 3 * time.Second

func (s *DroneAPI) OverrideMissionAndHover(ctx context.Context, command *common.MessageCommandLong) error {
	// 1. Guard clause to ensure the correct command is being executed
	if command.Command != common.MAV_CMD_DO_REPOSITION && command.Command != common.MAV_CMD_NAV_LOITER_UNLIM {
		return fmt.Errorf("invalid command ID %d: only loiter/reposition commands are allowed", command.Command)
	}

	log.Printf("Sending override command (%d): Hovering initiated...", command.Command)

	// 2. Send the command to the drone
	err := s.drone.node.WriteMessageAll(command)
	if err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}

	timer := time.NewTimer(commandTimeout)
	defer timer.Stop()

	// 3. Listen for the drone's acknowledgment response
	for {
		select {
		case <-ctx.Done():
			log.Println("Hover override command aborted by external signal.")
			return ctx.Err()

		case <-timer.C:
			log.Println("Timeout waiting for drone command acknowledgment.")
			return errors.New("command timed out without drone response")

		case evt, ok := <-s.drone.node.Events():
			if !ok {
				return errors.New("node events channel closed unexpectedly")
			}

			frm, ok := evt.(*gomavlib.EventFrame)
			if !ok {
				continue
			}

			// Filter: Only process messages coming from our targeted drone
			if frm.SystemID() != command.TargetSystem || frm.ComponentID() != command.TargetComponent {
				continue
			}

			switch msg := frm.Message().(type) {
			case *common.MessageCommandAck:
				// Ensure this ACK is for the command we just sent
				if msg.Command != command.Command {
					continue
				}

				switch msg.Result {
				case common.MAV_RESULT_ACCEPTED:
					log.Println("Hover command successfully ACCEPTED by the drone.")
					return nil
				case common.MAV_RESULT_TEMPORARILY_REJECTED:
					return errors.New("drone temporarily rejected command (is it armed/in a valid flight mode?)")
				case common.MAV_RESULT_DENIED:
					return errors.New("drone denied the hover command")
				case common.MAV_RESULT_UNSUPPORTED:
					return errors.New("drone does not support this hover command")
				default:
					return fmt.Errorf("command failed with drone error code: %d", msg.Result)
				}
			}
		}
	}
}
