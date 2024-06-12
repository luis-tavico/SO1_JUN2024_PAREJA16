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
	app.Post("/procesos/crear", crearProceso)
	app.Delete("/procesos/eliminar/:pid", eliminarProceso)
	app.Get("/procesos/:pid", buscarProceso)

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
	/*
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
			case 0: // Ejecución
				enEjecucion++
			case 1: // Suspendidos
				suspendidos++
			case 2: // Detenidos
				detenidos++
			case 3: // Zombies
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

		response := map[string]interface{}{
			"procesos": procesos,
			"info":     info,
		}

		return c.JSON(response)
	*/

	// Obtener los procesos del sistema
	processes, err := getProcesses()
	if err != nil {
		fmt.Println("Error al obtener los procesos del sistema: ", err)
		return c.Status(500).SendString("Error al obtener los procesos del sistema")

	}

	/*
		// Convertir la estructura a formato JSON
		json_processes, err := json.Marshal(processes)
		if err != nil {
			fmt.Println("Error al convertir los datos a JSON: ", err)
			return c.Status(500).SendString("Error al convertir los datos a JSON")
		}
	*/

	return c.JSON(processes)
}

func crearProceso(c *fiber.Ctx) error {
	cmd := exec.Command("sh", "-c", "sleep infinity &")
	err := cmd.Run()
	if err != nil {
		return c.Status(500).SendString("Error al crear el proceso")
	}

	return c.SendString("Proceso creado exitosamente")
}

func eliminarProceso(c *fiber.Ctx) error {
	pid := c.Params("pid")
	cmd := exec.Command("sh", "-c", "kill "+pid)
	err := cmd.Run()
	if err != nil {
		return c.Status(500).SendString("Error al eliminar el proceso")
	}

	return c.SendString("Proceso eliminado exitosamente")
}

func buscarProceso(c *fiber.Ctx) error {
	pidStr := c.Params("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return c.Status(400).SendString("PID inválido")
	}

	processes, err := getProcesses()
	if err != nil {
		return c.Status(500).SendString("Error al obtener los procesos del sistema")
	}

	for _, process := range processes.Processes {
		if process.Pid == pid {
			return c.JSON(process)
		}
	}

	return c.Status(404).SendString("Proceso no encontrado")
}

func insertDB() {
	for range time.Tick(time.Second * 1) {

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
