package factory

import (
	"fmt"
	"sync"

	".main.go/assemblyspot"
	".main.go/vehicle"
)

const assemblySpots int = 5

type Factory struct {
	AssemblingSpots chan *assemblyspot.AssemblySpot
}

func New() *Factory {
	factory := &Factory{
		AssemblingSpots: make(chan *assemblyspot.AssemblySpot, assemblySpots),
	}

	totalAssemblySpots := 0

	for {
		factory.AssemblingSpots <- &assemblyspot.AssemblySpot{}

		totalAssemblySpots++

		if totalAssemblySpots >= assemblySpots {
			break
		}
	}

	return factory
}

//HINT: this function is currently not returning anything, make it return right away every single vehicle once assembled,
//(Do not wait for all of them to be assembled to return them all, send each one ready over to main)
func (f *Factory) StartAssemblingProcess(amountOfVehicles int) {
	vehicleList := f.generateVehicleLots(amountOfVehicles)

	c := make(chan int, assemblySpots)
	var wg sync.WaitGroup

	for i := 0; i < len(vehicleList); i++ {
		c <- i
		wg.Add(1)

		go func(v *vehicle.Car) {
			defer wg.Done()
			fmt.Println("Assembling vehicle...")

			idleSpot := <-f.AssemblingSpots
			idleSpot.SetVehicle(v)

			vehicle, err := idleSpot.AssembleVehicle()

			if err == nil {
				vehicle.TestingLog = f.testCar(vehicle)
				vehicle.AssembleLog = idleSpot.GetAssembledLogs()
				idleSpot.SetVehicle(nil)

				fmt.Printf("The testingLog from car with id %d, is: %s \n", vehicle.Id, vehicle.TestingLog)
				fmt.Printf("The AssembleLog from car with id %d, is: %s \n", vehicle.Id, vehicle.AssembleLog)

				f.AssemblingSpots <- idleSpot
			}

			<-c
		}(&vehicleList[i])
	}

	wg.Wait()
}

func (Factory) generateVehicleLots(amountOfVehicles int) []vehicle.Car {
	var vehicles = []vehicle.Car{}
	var index = 0

	for {
		vehicles = append(vehicles, vehicle.Car{
			Id:            index,
			Chassis:       "NotSet",
			Tires:         "NotSet",
			Engine:        "NotSet",
			Electronics:   "NotSet",
			Dash:          "NotSet",
			Sits:          "NotSet",
			Windows:       "NotSet",
			EngineStarted: false,
		})

		index++

		if index >= amountOfVehicles {
			break
		}
	}

	return vehicles
}

func (f *Factory) testCar(car *vehicle.Car) string {
	logs := ""

	log, err := car.StartEngine()
	if err == nil {
		logs += log + ", "
	} else {
		logs += err.Error() + ", "
	}

	log, err = car.MoveForwards(10)
	if err == nil {
		logs += log + ", "
	} else {
		logs += err.Error() + ", "
	}

	log, err = car.MoveForwards(10)
	if err == nil {
		logs += log + ", "
	} else {
		logs += err.Error() + ", "
	}

	log, err = car.TurnLeft()
	if err == nil {
		logs += log + ", "
	} else {
		logs += err.Error() + ", "
	}

	log, err = car.TurnRight()
	if err == nil {
		logs += log + ", "
	} else {
		logs += err.Error() + ", "
	}

	log, err = car.StopEngine()
	if err == nil {
		logs += log + ", "
	} else {
		logs += err.Error() + ", "
	}

	return logs
}
