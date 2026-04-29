package gomavlinkdroneapi

import (
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
func (s *DroneAPI) UploadMission(missionItems *[]ardupilotmega.MessageMissionItemInt, targetSystem uint8, targetComponent uint8) bool {
	var rtn bool
	totalItems := uint16(len(*missionItems))

	// 1. Start the upload by sending MISSION_COUNT
	countMsg := &common.MessageMissionCount{
		TargetSystem:    targetSystem,
		TargetComponent: targetComponent,
		Count:           totalItems,
		MissionType:     common.MAV_MISSION_TYPE_MISSION,
	}
	s.drone.node.WriteMessageAll(countMsg)
	log.Printf("Sent MISSION_COUNT: %d waypoints\n", totalItems)

	// 2. Initialize a timeout timer (e.g., 10 seconds)
	timeoutDuration := 10 * time.Second
	timer := time.NewTimer(timeoutDuration)
	defer timer.Stop()

	// 3. Create a labeled loop to control breaks across the select statement
uploadLoop:
	for {
		select {
		// Handle incoming events from gomavlib
		case evt, ok := <-s.drone.node.Events():
			if !ok {
				log.Println("Event channel closed.")
				break uploadLoop
			}

			if frm, ok := evt.(*gomavlib.EventFrame); ok {
				switch msg := frm.Message().(type) {

				case *common.MessageMissionRequestInt:
					// Reset the timer on activity to allow long missions to complete
					if !timer.Stop() {
						<-timer.C
					}
					timer.Reset(timeoutDuration)

					if msg.Seq < totalItems {
						item := (*missionItems)[msg.Seq]
						s.drone.node.WriteMessageAll(&item)
						fmt.Printf("Uploaded waypoint #%d\n", msg.Seq)
					}

				case *common.MessageMissionRequest:
					if !timer.Stop() {
						<-timer.C
					}
					timer.Reset(timeoutDuration)

					seq := msg.Seq
					if int(seq) < len(*missionItems) {
						log.Printf("Uploading item %d...\n", seq)

						itemToSend := (*missionItems)[seq]
						itemToSend.TargetSystem = targetSystem
						itemToSend.TargetComponent = targetComponent

						s.drone.node.WriteMessageTo(frm.Channel, &itemToSend)
					}

				case *common.MessageMissionAck:
					if msg.Type == common.MAV_MISSION_ACCEPTED {
						fmt.Println("Mission uploaded and ACCEPTED by the drone!")
						rtn = true
					} else {
						log.Printf("Mission rejected by drone! Code: %d", msg.Type)
					}
					// Break the for loop because the handshake is complete
					break uploadLoop
				}
			}

		// Trigger if no successful activity happens within the window
		case <-timer.C:
			log.Println("Timeout reached: Drone did not complete mission handshake.")
			break uploadLoop
		}
	}
	return rtn
}

func (s *DroneAPI) DownloadMissions(targetSystem uint8, targetComponent uint8) (map[uint16]*ardupilotmega.MessageMissionItemInt, error) {

	var totalItems uint16
	var currentRequested uint16
	missionItems := make(map[uint16]*ardupilotmega.MessageMissionItemInt)

	log.Println("Requesting mission list from drone...")
	s.drone.node.WriteMessageAll(
		&ardupilotmega.MessageMissionRequestList{
			TargetSystem:    targetSystem,
			TargetComponent: targetComponent,
		},
	)

	// 1. Create a timer that resets whenever we get data
	timeoutDuration := 10 * time.Second
	timer := time.NewTimer(timeoutDuration)
	defer timer.Stop()

	// 2. Continuous loop with a select statement
	for {
		select {
		// Branch A: We received an event from gomavlib
		case evt, ok := <-s.drone.node.Events():
			if !ok {
				log.Println("Node events channel closed.")
				return nil, nil
			}

			switch e := evt.(type) {
			case *gomavlib.EventFrame:
				switch msg := e.Message().(type) {

				case *ardupilotmega.MessageMissionCount:
					// Data received! Reset the timeout timer
					if !timer.Stop() {
						<-timer.C
					}
					timer.Reset(timeoutDuration)

					totalItems = msg.Count
					log.Printf("Drone reported %d items. Starting download...", totalItems)

					if totalItems == 0 {
						log.Println("No mission items to download.")
						return nil, nil
					}

					currentRequested = 0
					s.drone.node.WriteMessageTo(e.Channel, &ardupilotmega.MessageMissionRequestInt{
						TargetSystem:    targetSystem,
						TargetComponent: targetComponent,
						Seq:             currentRequested,
					})

				case *ardupilotmega.MessageMissionItemInt:
					// Data received! Reset the timeout timer
					if !timer.Stop() {
						<-timer.C
					}
					timer.Reset(timeoutDuration)

					log.Printf("Received item #%d", msg.Seq)
					missionItems[msg.Seq] = msg

					if uint16(len(missionItems)) == totalItems {
						log.Println("All items received successfully!")

						s.drone.node.WriteMessageTo(e.Channel, &ardupilotmega.MessageMissionAck{
							TargetSystem:    targetSystem,
							TargetComponent: targetComponent,
							Type:            ardupilotmega.MAV_MISSION_ACCEPTED,
						})
						return missionItems, nil
					}

					currentRequested++
					s.drone.node.WriteMessageTo(e.Channel, &ardupilotmega.MessageMissionRequestInt{
						TargetSystem:    targetSystem,
						TargetComponent: targetComponent,
						Seq:             currentRequested,
					})
				}
			}

		// Branch B: No data was received before the timer fired
		case <-timer.C:
			log.Printf("Timeout! Drone failed to respond within %v. Aborting.", timeoutDuration)
			return nil, errors.New("Timeed out before the missions could be downloaded ")
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
func (s *DroneAPI) OverrideMissionAndHover(command *common.MessageCommandLong) error {
	err := s.drone.node.WriteMessageAll(command)
	log.Println("Sent override command: Hovering initiated.")
	return err
}
