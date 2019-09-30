package main

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/natefinch/lumberjack"
	rpio "github.com/stianeikeland/go-rpio"
	"github.com/tkanos/gonfig"
)

type Configuration struct {
	Name         *string
	Manufacturer *string
	Model        *string
	SerialNumber *string

	TogglePin int `json:"pin"`
}

var (
	projectRoot   string
	pin           *rpio.Pin
	configuration = Configuration{}
	openFilename  = "opened"
)

func saveState(isOpen bool) {
	if isOpen {
		file, err := os.Create(projectRoot + openFilename)
		if err != nil {
			log.Fatalln(err)
			return
		}
		defer file.Close()
	} else {
		err := os.Remove(projectRoot + openFilename)
		if err != nil {
			log.Fatalln(err)
			return
		}
	}
}

func getDoorState() int {
	if _, err := os.Stat(projectRoot + openFilename); os.IsNotExist(err) {
		log.Println("Door is Closed")
		return characteristic.CurrentDoorStateClosed
	}
	log.Println("Door is Open")
	return characteristic.CurrentDoorStateOpen
}

func toggleGPIO() {
	if pin != nil {
		pin.Output()
		pin.High()
		time.Sleep(time.Second)
		pin.Low()
	}
}

func accessoryInfo(config Configuration) accessory.Info {
	info := accessory.Info{
		Name:         "Garage Door",
		Manufacturer: "Liftmaster",
		Model:        "Professional 1/3 HP",
		SerialNumber: "l1f7m4573r",
	}
	if config.Name != nil {
		info.Name = *config.Name
	}
	if config.Manufacturer != nil {
		info.Manufacturer = *config.Manufacturer
	}
	if config.Model != nil {
		info.Model = *config.Model
	}
	if config.SerialNumber != nil {
		info.SerialNumber = *config.SerialNumber
	}
	return info
}

func updateDoorState(newState int, service *service.GarageDoorOpener) {
	log.Println("New State", newState)
	toggleGPIO()
	switch newState {
	case characteristic.TargetDoorStateOpen:
		log.Println("Opening")
		saveState(true)
		service.CurrentDoorState.SetValue(characteristic.CurrentDoorStateOpen)
		log.Println("Open")
	case characteristic.TargetDoorStateClosed:
		log.Println("Closing")
		saveState(false)
		service.CurrentDoorState.SetValue(characteristic.CurrentDoorStateClosed)
		log.Println("Closed")
	default:
		log.Fatalln("Error: Unknown State")
	}
	log.Println("Finished")
}

func setupGarageOpener() *accessory.Accessory {
	info := accessoryInfo(configuration)
	acc := accessory.New(info, accessory.TypeGarageDoorOpener)
	garageService := service.NewGarageDoorOpener()

	garageService.TargetDoorState.OnValueRemoteGet(func() int {
		log.Println("Getting Target Door State - Remotely")
		return getDoorState()
	})

	garageService.TargetDoorState.OnValueRemoteUpdate(func(newState int) {
		log.Println("Updating Target Door State - Remotely")
		updateDoorState(newState, garageService)
	})

	garageService.CurrentDoorState.OnValueRemoteGet(func() int {
		log.Println("Getting Current Door State - Remotely")
		return getDoorState()
	})
	garageService.CurrentDoorState.OnValueGet(func() interface{} {
		log.Println("Getting Current Door State")
		return getDoorState()
	})
	// Initial value should be closed
	garageService.TargetDoorState.SetValue(characteristic.TargetDoorStateClosed)
	garageService.CurrentDoorState.SetValue(characteristic.CurrentDoorStateClosed)
	acc.AddService(garageService.Service)
	return acc
}

func main() {
	e, err := os.Executable()
	if err != nil {
		panic(err)
	}
	projectRoot = path.Dir(e) + "/"
	// Setup Rotating Log
	log.SetOutput(&lumberjack.Logger{
		Filename:   projectRoot + "log",
		MaxSize:    500, // megabytes
		MaxBackups: 1,
		MaxAge:     28, //days
	})
	if err := gonfig.GetConf(projectRoot+"config.json", &configuration); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	log.Println("Loaded Configuration")
	// Setup GPIO
	if err := rpio.Open(); err != nil {
		log.Println(err)
	} else {
		log.Println("Pin=", configuration.TogglePin)
		togglePin := rpio.Pin(configuration.TogglePin)
		pin = &togglePin
		pin.Output()
		pin.Low()
	}
	// Unmap gpio memory when done
	defer rpio.Close()
	acc := setupGarageOpener()
	log.Println("Created Garage Accessory")
	t, err := hc.NewIPTransport(hc.Config{StoragePath: projectRoot, Pin: "32191123"}, acc)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Println("Created IP Transport")
	hc.OnTermination(func() {
		log.Println("Stopping")
		t.Stop()
		os.Exit(1)
	})
	log.Println("Starting IP Transport")
	t.Start()
}
