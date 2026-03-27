package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/grid-x/modbus"
)

const (
	// Holding Register Addresses (Read/Write) - untuk data utama
	LeavingChilledWaterTempSettings = 40001 // Setting temperature (target)
	EnteringChilledWaterTemp        = 40002 // Temperature masuk
	LeavingChilledWaterTemp         = 40003 // Temperature keluar

	// Coil Addresses (Read/Write) - untuk kontrol
	ControlEnablePublish = 300 // Enable/disable publishing
	ControlResetValues   = 301 // Reset ke default values
	ControlSimulation    = 302 // Enable/disable simulasi

	// Modbus Server Config
	ModbusIP   = "localhost"
	ModbusPort = 503
)

type ChillerData struct {
	LeavingTempSettings float32 // Setting temperature
	EnteringTemp        float32 // Entering water temp
	LeavingTemp         float32 // Leaving water temp
	PublishEnabled      bool
	SimulationEnabled   bool
}

func main() {
	// Setup Modbus TCP Server
	handler := modbus.TCPClient(fmt.Sprintf("%s:%d", ModbusIP, ModbusPort))
	handler.ReadHoldingRegisters() = make([]byte, 65536)
	handler.Coils = make([]byte, 8192)

	server := modbus.NewTCPServer(handler)

	log.Printf("Starting Modbus TCP Server on %s:%d", ModbusIP, ModbusPort)

	// Start server in goroutine
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Initialize default values
	chiller := ChillerData{
		LeavingTempSettings: 7.0,  // Default setting 7°C
		EnteringTemp:        12.0, // Default entering temp
		LeavingTemp:         7.5,  // Default leaving temp
		PublishEnabled:      true,
		SimulationEnabled:   true,
	}

	// Write initial control values
	setCoil(handler, ControlEnablePublish, true)
	setCoil(handler, ControlSimulation, true)
	setCoil(handler, ControlResetValues, false)

	log.Println("Server started successfully!")
	log.Println("\nMemory Locations:")
	log.Printf("  40001: Leaving Chilled Water Temp Settings")
	log.Printf("  40002: Entering Chilled Water Temp")
	log.Printf("  40003: Leaving Chilled Water Temp")
	log.Println("\nControl Coils:")
	log.Printf("  300: Enable/Disable Publishing")
	log.Printf("  301: Reset to Default Values")
	log.Printf("  302: Enable/Disable Simulation")
	log.Println("\nServer is running. Press Ctrl+C to stop.")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Read control coils
		chiller.PublishEnabled = getCoil(handler, ControlEnablePublish)
		chiller.SimulationEnabled = getCoil(handler, ControlSimulation)
		resetRequested := getCoil(handler, ControlResetValues)

		// Handle reset
		if resetRequested {
			log.Println("Reset requested - restoring default values")
			chiller.LeavingTempSettings = 7.0
			chiller.EnteringTemp = 12.0
			chiller.LeavingTemp = 7.5
			setCoil(handler, ControlResetValues, false)
		}

		// Simulate temperature changes if enabled
		if chiller.SimulationEnabled {
			// Read current setting from register (bisa diubah dari luar)
			currentSetting := getFloat32(handler, LeavingChilledWaterTempSettings)
			if currentSetting > 0 && currentSetting < 30 {
				chiller.LeavingTempSettings = currentSetting
			}

			// Simulate entering temperature (varies slightly)
			chiller.EnteringTemp = 12.0 + float32(math.Sin(float64(time.Now().Unix())/10.0)*2.0)

			// Simulate leaving temperature (tries to reach setting)
			diff := chiller.LeavingTempSettings - chiller.LeavingTemp
			chiller.LeavingTemp += diff * 0.1 // Gradual approach to target
		}

		// Publish data if enabled
		if chiller.PublishEnabled {
			setFloat32(handler, LeavingChilledWaterTempSettings, chiller.LeavingTempSettings)
			setFloat32(handler, EnteringChilledWaterTemp, chiller.EnteringTemp)
			setFloat32(handler, LeavingChilledWaterTemp, chiller.LeavingTemp)

			log.Printf("Published - Setting: %.2f°C | Entering: %.2f°C | Leaving: %.2f°C",
				chiller.LeavingTempSettings,
				chiller.EnteringTemp,
				chiller.LeavingTemp)
		} else {
			log.Println("Publishing disabled by control coil 300")
		}
	}
}

// Helper functions untuk float32
func setFloat32(handler *modbus.TCPServerHandler, address uint16, value float32) {
	bits := math.Float32bits(value)
	// Modbus address dimulai dari 1, tapi array index dari 0
	index := (address - 1) * 2
	binary.BigEndian.PutUint16(handler.HoldingRegisters[index:], uint16(bits>>16))
	binary.BigEndian.PutUint16(handler.HoldingRegisters[index+2:], uint16(bits))
}

func getFloat32(handler *modbus.TCPServerHandler, address uint16) float32 {
	index := (address - 1) * 2
	high := binary.BigEndian.Uint16(handler.HoldingRegisters[index:])
	low := binary.BigEndian.Uint16(handler.HoldingRegisters[index+2:])
	bits := uint32(high)<<16 | uint32(low)
	return math.Float32frombits(bits)
}

// Helper functions untuk coils
func setCoil(handler *modbus.TCPServerHandler, address uint16, value bool) {
	byteIndex := address / 8
	bitIndex := address % 8
	if value {
		handler.Coils[byteIndex] |= 1 << bitIndex
	} else {
		handler.Coils[byteIndex] &^= 1 << bitIndex
	}
}

func getCoil(handler *modbus.TCPServerHandler, address uint16) bool {
	byteIndex := address / 8
	bitIndex := address % 8
	return (handler.Coils[byteIndex] & (1 << bitIndex)) != 0
}
