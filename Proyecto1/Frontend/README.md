# Frontend

## App

Este es el componente principal de la aplicación que maneja la navegación entre diferentes componentes según la ruta actual.

### Imports

Se importan las bibliotecas y componentes necesarios para el funcionamiento de la aplicación:

- `useState` de React.
- Logos de React, Bootstrap y Vite.
- Componentes personalizados `ModoOscuro`, `Head`, `Estadisticas` y `TablaProcesos`.
- Estilos personalizados desde `./styles/App.css`.

### Estados del Componente

No se utilizan estados locales en este componente.

### Lógica de Rutas

Se utiliza un `switch` para determinar qué componente renderizar basado en la ruta actual (`window.location.pathname`):

- `/` y `/estadisticas` renderizan el componente `Estadisticas`.
- `/tablaprocesos` renderiza el componente `TablaProcesos`.
- Otras rutas pueden añadirse en casos adicionales dentro del `switch`.

### Renderizado del Componente

El componente principal renderiza el componente `Head` y el componente correspondiente a la ruta actual.


## Barra de Navegación

Este componente de React crea una barra de navegación utilizando los componentes `Nav` y `Navbar` de `react-bootstrap`. La barra de navegación incluye enlaces a las secciones de estadísticas y tabla de procesos.

### Imports

Se importan los componentes necesarios desde `react-bootstrap`.

### Renderizado del Componente

El componente `Head` retorna una barra de navegación (`Navbar`) con un tema oscuro, que se expande en pantallas grandes. Incluye un enlace de marca a la ruta `/estadisticas` y dos enlaces de navegación a las rutas `/estadisticas` y `/tablaprocesos`.

## Estadísticas

Este componente de React muestra estadísticas en gráficos de tipo donut (doughnut) usando la biblioteca `react-chartjs-2` para visualizar el porcentaje de uso de CPU y RAM en tiempo real.

### Imports

Se importan las bibliotecas necesarias para el funcionamiento del componente, incluyendo React, hooks (`useState`, `useEffect`), `react-chartjs-2` para los gráficos y los estilos personalizados desde `../styles/Estilo.css`.

### Estados del Componente

- `cpuData`: Almacena el porcentaje de uso de la CPU obtenido de la API.
- `ramData`: Almacena el porcentaje de uso de la RAM obtenido de la API.

### useEffect

Se utiliza para cargar los datos de las estadísticas desde una API al montar el componente y actualizarlos cada 2 segundos.

### Función `fetchData`

Realiza una llamada a la API para obtener los datos de uso de CPU y RAM, y actualiza los estados `cpuData` y `ramData`.

### Función `doughnutData`

Recibe una etiqueta y un porcentaje de uso, y devuelve los datos formateados para el gráfico de tipo donut, incluyendo los colores y etiquetas correspondientes.

### Renderizado del Componente

Incluye un contenedor principal con el título "SO1 - JUN 2024" y dos gráficos de tipo donut que muestran el porcentaje de uso de CPU y RAM. Los gráficos se renderizan solo si los datos están disponibles.

### Estilo

Asegúrate de crear el archivo `../styles/Estilo.css` para aplicar los estilos necesarios al componente.



## Tabla de Procesos

Este componente de React muestra una tabla de procesos, permitiendo la creación y eliminación de procesos, así como la búsqueda por nombre o PID.

### Imports

Se importan las bibliotecas necesarias para el funcionamiento del componente, incluyendo React, hooks (`useState`, `useEffect`), funciones de `react-table` y componentes de `react-bootstrap`.

### Estados del Componente

- `processes`: Almacena la lista de procesos obtenidos de la API.
- `searchTerm`: Término de búsqueda por nombre.
- `searchPidTerm`: Término de búsqueda por PID.
- `info`: Información adicional sobre los procesos.
- `pidToKill`: PID del proceso que se desea eliminar.
- `searchBy`: Determina si se busca por nombre o por PID.

### useEffect

Se utiliza para cargar los datos de los procesos desde una API al montar el componente.

### Columnas de la Tabla

Definidas usando `React.useMemo` para mejorar el rendimiento, describen cómo se deben mostrar los datos en la tabla.

### Funciones de Manejo

- `handleCreateProcess`: Para crear un nuevo proceso.
- `handleKillProcess`: Para eliminar un proceso usando el PID almacenado.
- `handleSearchChange` y `handleSearchPidChange`: Para manejar los cambios en los campos de búsqueda.
- `handleSearchByChange`: Para cambiar entre búsqueda por nombre y búsqueda por PID.

### Filtrado de Filas

Se filtran las filas de la tabla según el término de búsqueda y el tipo de búsqueda seleccionado.

### Renderizado del Componente

Incluye botones para crear y eliminar procesos, campos de entrada para buscar procesos, y una tabla que muestra los datos filtrados de los procesos. La tabla puede expandirse para mostrar información adicional sobre los procesos hijos si los hay.
