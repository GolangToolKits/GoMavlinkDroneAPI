package gomavlinkdroneapi_test

import (
	"context"
	"testing"

	gomavlinkdroneapi "github.com/GolangToolKits/GoMavlinkDroneAPI"
	"github.com/bluenviron/gomavlib/v3"
)

func TestDroneAPI_ClearGeofence(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		targetSystem    uint8
		targetComponent uint8
		want            bool
		clientAddress   string
		outVersion      gomavlib.Version
		outSystemID     byte
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
			got := s.ClearGeofence(ctx, tt.targetSystem, tt.targetComponent)
			// TODO: update the condition below to compare got with tt.want.
			if got != false {
				t.Errorf("ClearGeofence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDroneAPI_UploadGeofence(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		newFence        []gomavlinkdroneapi.GeoFence
		targetSystem    uint8
		targetComponent uint8
		want            bool
		wantErr         bool
		clientAddress   string
		outVersion      gomavlib.Version
		outSystemID     byte
	}{
		// TODO: Add test cases.
		// {
		// 	name:          "test 1",
		// 	clientAddress: "1.2.3.4:5600",
		// 	outVersion:    gomavlib.V2,
		// 	outSystemID:   10,
		// 	wantErr:       false,
		// 	newFence: &[]gomavlinkdroneapi.GeoFence{
		// 		{Lat: 473977418, Lon: 85455939},
		// 		{Lat: 473977418, Lon: 85465939},
		// 		{Lat: 473987418, Lon: 85465939},
		// 		{Lat: 473987418, Lon: 85455939},
		// 	},
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
			got, gotErr := s.UploadGeofence(ctx, tt.newFence, tt.targetSystem, tt.targetComponent)
			if gotErr == nil {
				if !tt.wantErr {
					t.Errorf("UploadGeofence() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("UploadGeofence() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if got {
				t.Errorf("UploadGeofence() = %v, want %v", got, tt.want)
			}
		})
	}
}
