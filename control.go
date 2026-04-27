package gomavlinkdroneapi

import (
	"log"
	"time"

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
