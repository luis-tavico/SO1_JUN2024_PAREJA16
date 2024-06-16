package main

import (
	"Backend/Controller"
	"Backend/Database"
	"Backend/Model"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// Función de inicialización
func initSystem() error {
	// Verificar si /proc/ram_so1_jun2024 existe
	cmd := exec.Command("sh", "-c", "cat /proc/ram_so1_jun2024")
	_, err := cmd.CombinedOutput()
	if err == nil {
		// Verificar si /proc/cpu_so1_1s2024 existe
		cmd = exec.Command("sh", "-c", "cat /proc/cpu_so1_1s2024")
		_, err = cmd.CombinedOutput()
		if err == nil {
			// Ambos archivos existen, continuar con el flujo normal
			return nil
		}
	}

	// Si alguno de los archivos no existe, proceder con la compilación e inserción de módulos
	baseDir := "/home/oscar/SO1_JUN2024_PAREJA16/Proyecto1/Modules" // Cambia esta ruta según tu estructura de directorios

	// Verificar y cargar módulo de CPU
	if err := checkAndLoadModule(filepath.Join(baseDir, "CPU"), "cpu.ko"); err != nil {
		return fmt.Errorf("error al cargar el módulo de CPU: %w", err)
	}

	// Verificar y cargar módulo de RAM
	if err := checkAndLoadModule(filepath.Join(baseDir, "RAM"), "ram.ko"); err != nil {
		return fmt.Errorf("error al cargar el módulo de RAM: %w", err)
	}

	return nil
}

// Función para verificar y cargar módulos
func checkAndLoadModule(dir, moduleName string) error {
	modulePath := filepath.Join(dir, moduleName)
	if _, err := os.Stat(modulePath); os.IsNotExist(err) {
		// El archivo no existe, ejecutar make
		cmd := exec.Command("make")
		cmd.Dir = dir
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("error al ejecutar make en %s: %s", dir, output)
		}
	}

	// Cargar el módulo
	cmd := exec.Command("sudo", "insmod", modulePath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error al cargar el módulo %s: %s", modulePath, output)
	}

	return nil
}

func main() {

	// Inicializar el sistema
	if err := initSystem(); err != nil {
		log.Fatalf("Error al inicializar el sistema: %v", err)
	}
	app := fiber.New()

	// Habilitar CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,DELETE",
		AllowHeaders: "Content-Type, Authorization",
	}))

	if err := Database.Connect(); err != nil {
		log.Fatal("Error en", err)
	}

	// Rutas API
	app.Get("/estadisticas", getEstadisticas)
	app.Get("/procesos", getProcesos)
	app.Get("/procesos/crear", crearProceso)
	app.Post("/procesos/eliminar/:pid", detenerProceso)

	go insertDB()

	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}

}

// Estructura para los datos del sistema
type SystemData struct {
	RAM_percentage int     `json:"ram_percentage"`
	CPU_percentage float64 `json:"cpu_percentage"`
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
func getCPUdata() (float64, error) {
    cmd := exec.Command("sh", "-c", "mpstat 1 1 | awk '/Average:/ {print $12}'")
    stdout, err := cmd.CombinedOutput()
    if err != nil {
        return 0, err
    }

    // Convertir la salida a float64
    cpuIdleStr := strings.TrimSpace(string(stdout))
    cpuIdle, err := strconv.ParseFloat(cpuIdleStr, 64)
    if err != nil {
        return 0, err
    }

    // Calcular el porcentaje de uso de CPU
    cpuUsage := 100 - cpuIdle

	// Formatear a dos decimales
    cpuUsageStr := fmt.Sprintf("%.2f", cpuUsage)
    cpuUsage, err = strconv.ParseFloat(cpuUsageStr, 64)
    if err != nil {
        return 0, err
    }
	
    return cpuUsage, nil
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
		return c.Status(500).SendString(fmt.Sprintf("Error al obtener datos de la RAM: %v", err))
	}

	// Obtener datos de la CPU
	usedCPUPercentage, err := getCPUdata()
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("Error al obtener datos de la CPU: %v", err))
	}

	usedRAMPercentage := 100 - freeRAMPercentage

	estadisticas := map[string]interface{}{
		"ram_percentage": usedRAMPercentage,
		"cpu_percentage": usedCPUPercentage,
	}

	return c.JSON(estadisticas)
}

func getProcesos(c *fiber.Ctx) error {

	processData, err := getProcesses()
	if err != nil {
		return c.Status(500).SendString("Error al obtener los procesos del sistema aaa")
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

	// Actualizar la tabla de procesos
	err = updateProcessTable()
	if err != nil {
		response := map[string]string{
			"estado": "Error",
			"pid":    "Error al actualizar la tabla de procesos",
		}
		return c.JSON(response)
	}

	// Crear la respuesta
	response := map[string]string{
		"estado": "creado",
		"pid":    pid,
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

	// Actualizar la tabla de procesos
	err = updateProcessTable()
	if err != nil {
		response := map[string]string{
			"estado": "Error",
			"pid":    "Error al actualizar la tabla de procesos",
		}
		return c.JSON(response)
	}

	// Crear la respuesta
	response := map[string]string{
		"estado": "Eliminado",
		"pid":    "El proceso " + pid + " fue detenido.",
	}

	// Retornar la respuesta JSON
	return c.JSON(response)
}

func updateProcessTable() error {
	processes, err := getProcesses()
	if err != nil {
		return err
	}

	// Primero elimina los datos antiguos
	Controller.DeleteDataProcess("process")

	// Luego inserta los datos actualizados
	for _, process := range processes.Processes {
		dataProcess := Model.DataProcess{
			Pid:   strconv.Itoa(process.Pid),
			Name:  process.Name,
			User:  strconv.Itoa(process.User),
			State: strconv.Itoa(process.State),
			Ram:   strconv.Itoa(process.Ram),
		}
		Controller.InsertDataProcess("process", dataProcess)

		for _, child := range process.Child {
			dataProcess := Model.DataProcess{
				Pid:      strconv.Itoa(child.Pid),
				Name:     child.Name,
				State:    strconv.Itoa(child.State),
				PidPadre: strconv.Itoa(process.Pid),
			}
			Controller.InsertDataProcess("process", dataProcess)
		}
	}

	return nil
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
			Used_percentage: strconv.FormatFloat(usedCPUPercentage, 'f', 2, 64),
			Free_percentage: strconv.FormatFloat(freeCPUPercentage, 'f', 2, 64),
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