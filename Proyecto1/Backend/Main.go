package main

import (
	"Backend/Controller"
	"Backend/Database"
	"Backend/Model"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// Habilitar CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5174",
		AllowMethods: "GET,POST,DELETE",
	}))

	if err := Database.Connect(); err != nil {
		log.Fatal("Error en", err)
	}

	// Rutas API
	app.Get("/estadisticas", getEstadisticas)
	app.Get("/procesos", getProcesos)
	app.Get("/procesos/crear", crearProceso)
	app.Post("/procesos/eliminar/:pid", detenerProceso)

	//go insertDB()

	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}

}

// Estructura para los datos del sistema
type SystemData struct {
	RAM_percentage int `json:"ram_percentage"`
	CPU_percentage int `json:"cpu_percentage"`
}

// Estructura para los procesos hijo
type ProcesoHijo struct {
	Pid      int    `json:"pid"`
	Name     string `json:"name"`
	State    int    `json:"state"`
	PidPadre int    `json:"pidPadre"`
}

// Estructura para los procesos padre
type ProcesoPadre struct {
	Pid   int           `json:"pid"`
	Name  string        `json:"name"`
	User  int           `json:"user"`
	State int           `json:"state"`
	Ram   int           `json:"ram"`
	Child []ProcesoHijo `json:"child"`
}

// Definicion de estructura para el objeto JSON
type ProcessData struct {
	Processes []ProcesoPadre `json:"processes"`
}

// Funcion para obtener datos de la RAM
func getRAMdata() (int, error) {
	cmd := exec.Command("sh", "-c", "cat /proc/ram_so1_jun2024")
	stdout, err := cmd.CombinedOutput()

	if err != nil {
		return 0, err
	}

	// Convertir la salida a formato JSON
	var data SystemData
	err = json.Unmarshal(stdout, &data)
	if err != nil {
		return 0, err
	}

	return data.RAM_percentage, nil
}

// Funcion para obtener datos de la CPU
func getCPUdata() (int, error) {
	cmd := exec.Command("sh", "-c", "cat /proc/cpu_so1_1s2024")
	stdout, err := cmd.CombinedOutput()

	if err != nil {
		return 0, err
	}

	// Convertir la salida a formato JSON
	var data SystemData
	err = json.Unmarshal(stdout, &data)
	if err != nil {
		return 0, err
	}

	return data.CPU_percentage, nil
}

// Funcion para obtener los procesos del sistema
func getProcesses() (ProcessData, error) {
	cmd := exec.Command("sh", "-c", "cat /proc/cpu_so1_1s2024")
	stdout, err := cmd.CombinedOutput()

	if err != nil {
		return ProcessData{}, err
	}
	// Convertir la salida a formato JSON
	var data ProcessData
	err = json.Unmarshal(stdout, &data)
	if err != nil {
		return ProcessData{}, err
	}

	return data, nil
}

func getEstadisticas(c *fiber.Ctx) error {
	// Obtener datos de la RAM
	freeRAMPercentage, err := getRAMdata()
	if err != nil {
		return c.Status(500).SendString("Error al obtener datos de la RAM")
	}

	// Obtener datos de la CPU
	usedCPUPercentage, err := getCPUdata()
	if err != nil {
		return c.Status(500).SendString("Error al obtener datos de la CPU")
	}

	usedRAMPercentage := 100 - freeRAMPercentage
	//freeCPUPercentage := 100 - usedCPUPercentage

	estadisticas := map[string]int{
		"ram_percentage": usedRAMPercentage,
		"cpu_percentage": usedCPUPercentage,
	}

	return c.JSON(estadisticas)
}

func getProcesos(c *fiber.Ctx) error {
	
	processData, err := getProcesses()
	if err != nil {
		return c.Status(500).SendString("Error al obtener los procesos del sistema")
	}

	var procesos []map[string]interface{}
	var enEjecucion, suspendidos, detenidos, zombies int

	for _, process := range processData.Processes {
		proceso := map[string]interface{}{
			"pid":     process.Pid,
			"nombre":  process.Name,
			"estado":  process.State,
			"ram":     process.Ram,
			"usuario": process.User,
		}
		procesos = append(procesos, proceso)

		// Contar los procesos por estado
		switch process.State {
		case 1: // Ejecución
			enEjecucion++
		case 0: // Suspendidos
			suspendidos++
		case 128: // Detenidos
			detenidos++
		case 1026: // Zombies
			zombies++
		}
	}

	total := len(procesos)
	info := map[string]int{
		"en_ejecucion": enEjecucion,
		"suspendidos":  suspendidos,
		"detenidos":    detenidos,
		"zombies":      zombies,
		"total":        total,
	}

	

// Obtener los procesos del sistema
processes, err := getProcesses()
if err != nil {
	fmt.Println("Error al obtener los procesos del sistema: ", err)
	return c.Status(500).SendString("Error al obtener los procesos del sistema")

}


response := map[string]interface{}{
	"procesos": processes,
	"info":     info,
}

return c.JSON(response)

}

func crearProceso(c *fiber.Ctx) error {
    // Crear el comando
    cmd := exec.Command("sleep", "infinity")

    // Iniciar el comando
    err := cmd.Start()
    if err != nil {
        response := map[string]string{
            "estado": "Error",
            "pid":    "Error al crear el proceso",
        }
        return c.JSON(response)
    }

    // Obtener el PID del proceso y convertirlo a string
    pid := strconv.Itoa(cmd.Process.Pid)

    // Crear la respuesta
    response := map[string]string{
        "estado": "creado",
        "pid":    "El proceso "+pid+" fue creado.",
    }

    // Retornar la respuesta JSON
    return c.JSON(response)
}

func detenerProceso(c *fiber.Ctx) error {
	pid := c.Params("pid")
	cmd := exec.Command("kill", "-9", pid)
	err := cmd.Run()
	if err != nil {
		response := map[string]string{
			"estado": "Error",
			"pid":    "Error al detener el proceso",
		}
		return c.JSON(response)
	}
	// Crear la respuesta
    response := map[string]string{
        "estado": "Eliminado",
        "pid":    "El proceso "+pid+ " fue detenido.",
    }

    // Retornar la respuesta JSON
    return c.JSON(response)
}

func insertDB() {
	for range time.Tick(time.Second * 1) {

		// Eliminar todos los documentos de la colección ram
		Controller.DeleteDataRAM("ram")

		// Eliminar todos los documentos de la colección cpu
		Controller.DeleteDataCPU("cpu")

		// Eliminar todos los documentos de la colección process
		Controller.DeleteDataProcess("process")

		// RAM
		freeRAMPercentage, err := getRAMdata()
		if err != nil {
			fmt.Println("Error al obtener datos de la RAM: ", err)
		}

		usedRAMPercentage := 100 - freeRAMPercentage

		dataRAM := Model.DataRAM{
			Used_percentage: strconv.Itoa(usedRAMPercentage),
			Free_percentage: strconv.Itoa(freeRAMPercentage),
		}

		// Insertar los datos en la base de datos
		Controller.InsertDataRAM("ram", dataRAM)

		// CPU
		usedCPUPercentage, err := getCPUdata()
		if err != nil {
			fmt.Println("Error al obtener datos de la CPU: ", err)
		}

		freeCPUPercentage := 100 - usedCPUPercentage

		dataCPU := Model.DataCPU{
			Used_percentage: strconv.Itoa(usedCPUPercentage),
			Free_percentage: strconv.Itoa(freeCPUPercentage),
		}
		// Insertar los datos en la base de datos
		Controller.InsertDataCPU("cpu", dataCPU)

		// Procesos
		processes, err := getProcesses()
		if err != nil {
			fmt.Println("Error al obtener los procesos del sistema: ", err)
		}

		for _, process := range processes.Processes {
			dataProcess := Model.DataProcess{
				Pid:   strconv.Itoa(process.Pid),
				Name:  process.Name,
				User:  strconv.Itoa(process.User),
				State: strconv.Itoa(process.State),
				Ram:   strconv.Itoa(process.Ram),
			}
			// Insertar los datos en la base de datos
			Controller.InsertDataProcess("process", dataProcess)

			for _, child := range process.Child {
				dataProcess := Model.DataProcess{
					Pid:      strconv.Itoa(child.Pid),
					Name:     child.Name,
					State:    strconv.Itoa(child.State),
					PidPadre: strconv.Itoa(process.Pid),
				}
				// Insertar los datos en la base de datos
				Controller.InsertDataProcess("process", dataProcess)
			}
		}

	}
}
