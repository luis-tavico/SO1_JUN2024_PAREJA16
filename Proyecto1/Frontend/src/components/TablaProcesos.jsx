import React, { useState, useEffect } from 'react'; // Importa React y hooks de estado y efectos
import { useTable, useExpanded } from 'react-table'; // Importa funciones de react-table para crear tablas y manejar filas expandidas
import { Table, Button, InputGroup, FormControl } from 'react-bootstrap'; // Importa componentes de react-bootstrap para la interfaz
import '../components/Estilo.css'; // Importa estilos personalizados

function TablaProcesos() {
  // Define los estados del componente
  const [processes, setProcesses] = useState([]); // Estado para almacenar los procesos
  const [searchTerm, setSearchTerm] = useState(''); // Estado para el término de búsqueda por nombre
  const [searchPidTerm, setSearchPidTerm] = useState(''); // Estado para el término de búsqueda por PID
  const [info, setInfo] = useState({}); // Estado para la información adicional de los procesos
  const [pidToKill, setPidToKill] = useState(''); // Estado para almacenar el PID a eliminar
  const [searchBy, setSearchBy] = useState('name'); // Estado para determinar si se busca por nombre o por PID

  // Efecto para cargar los datos de los procesos desde la API
  useEffect(() => {
    fetch('http://192.168.122.195:8080/procesos')
      .then(response => response.json())
      .then(data => {
        if (data && data.processes) {
          setProcesses(data.processes); // Actualiza los procesos si se obtienen datos válidos
          setInfo(data.info || {}); // Actualiza la información adicional
        } else {
          setProcesses([]); // Si no hay datos válidos, limpia los procesos
          setInfo({}); // Limpia la información adicional
        }
      })
      .catch(error => {
        console.error('Error fetching processes:', error); // Muestra un error en caso de fallo
        setProcesses([]); // Limpia los procesos en caso de error
        setInfo({}); // Limpia la información adicional en caso de error
      });
  }, []);

  // Define las columnas de la tabla
  const columns = React.useMemo(
    () => [
      {
        Header: 'PID',
        accessor: 'pid', // Accede al campo 'pid' en los datos
      },
      {
        Header: 'Name',
        accessor: 'name', // Accede al campo 'name' en los datos
      },
      {
        Header: 'User',
        accessor: 'user', // Accede al campo 'user' en los datos
      },
      {
        Header: 'State',
        accessor: 'state', // Accede al campo 'state' en los datos
      },
      {
        Header: 'RAM',
        accessor: 'ram', // Accede al campo 'ram' en los datos
      },
      {
        Header: 'Actions',
        id: 'expander', // Columna para manejar la expansión de filas
        Cell: ({ row }) => (
          row.original.child && row.original.child.length > 0 ? (
            <span {...row.getToggleRowExpandedProps()}>
              {row.isExpanded ? '▼' : '▶'}
            </span>
          ) : null // No muestra ícono si no hay hijos
        ),
      },
    ],
    []
  );

  // Función para crear un proceso
  const handleCreateProcess = () => {
    fetch('/procesos/crear', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
    })
      .then(response => response.json())
      .then(data => {
        if (data && data.procesos) {
          setProcesses(data.procesos); // Actualiza los procesos con los nuevos datos
          setInfo(data.info || {}); // Actualiza la información adicional
        }
      })
      .catch(error => console.error('Error creating process:', error)); // Muestra un error en caso de fallo
  };

  // Función para eliminar un proceso
  const handleKillProcess = () => {
    fetch(`/procesos/eliminar/${pidToKill}`, {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
    })
      .then(response => response.json())
      .then(data => {
        if (data && data.procesos) {
          setProcesses(data.procesos); // Actualiza los procesos con los nuevos datos
          setInfo(data.info || {}); // Actualiza la información adicional
        } else {
          alert('PID no encontrado'); // Muestra una alerta si el PID no se encuentra
        }
      })
      .catch(error => console.error('Error killing process:', error)); // Muestra un error en caso de fallo
  };

  // Prepara los datos para la tabla
  const data = React.useMemo(
    () =>
      processes.map(process => ({
        ...process,
        subRows: process.child || [], // Agrega subfilas si hay procesos hijos
      })),
    [processes]
  );

  // Configura la tabla con react-table
  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    rows,
    prepareRow,
  } = useTable(
    {
      columns,
      data,
    },
    useExpanded // Habilita la expansión de filas
  );

  // Maneja el cambio en el término de búsqueda por nombre
  const handleSearchChange = event => {
    setSearchTerm(event.target.value);
  };

  // Maneja el cambio en el término de búsqueda por PID
  const handleSearchPidChange = event => {
    setSearchPidTerm(event.target.value);
  };

  // Cambia entre buscar por nombre o PID
  const handleSearchByChange = () => {
    setSearchBy(prevSearchBy => (prevSearchBy === 'name' ? 'pid' : 'name'));
  };

  // Filtra las filas según el término de búsqueda y el tipo de búsqueda
  const filteredRows = rows.filter(row => {
    if (searchBy === 'name') {
      return row.original.name.toLowerCase().includes(searchTerm.toLowerCase());
    }
    return row.original.pid.toString().includes(searchPidTerm);
  });

  return (
    <div>
      <div className="title-container">
        <h1>Tabla de Procesos</h1>
      </div>
      <Button onClick={handleCreateProcess}>Crear Proceso</Button>
      <InputGroup className="mb-3">
        <FormControl
          placeholder="PID para eliminar"
          value={pidToKill}
          onChange={(e) => setPidToKill(e.target.value)}
        />
        <Button variant="danger" onClick={handleKillProcess}>Eliminar Proceso</Button>
      </InputGroup>
      <Button onClick={handleSearchByChange}>
        {searchBy === 'name' ? 'Buscar por PID' : 'Buscar por Nombre'}
      </Button>
      {searchBy === 'name' ? (
        <input
          type="text"
          placeholder="Buscar por nombre"
          value={searchTerm}
          onChange={handleSearchChange}
        />
      ) : (
        <input
          type="text"
          placeholder="Buscar por PID"
          value={searchPidTerm}
          onChange={handleSearchPidChange}
        />
      )}
      <Table {...getTableProps()} striped bordered hover>
        <thead>
          {headerGroups.map(headerGroup => (
            <tr {...headerGroup.getHeaderGroupProps()}>
              {headerGroup.headers.map(column => (
                <th {...column.getHeaderProps()}>{column.render('Header')}</th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody {...getTableBodyProps()}>
          {filteredRows.map(row => {
            prepareRow(row);
            return (
              <React.Fragment key={row.id}>
                <tr {...row.getRowProps()}>
                  {row.cells.map(cell => (
                    <td {...cell.getCellProps()}>{cell.render('Cell')}</td>
                  ))}
                </tr>
                {row.isExpanded && row.original.child.length > 0 && (
                  <tr>
                    <td colSpan={columns.length}>
                      <Table striped bordered hover>
                        <thead>
                          <tr>
                            <th>PID</th>
                            <th>Name</th>
                            <th>State</th>
                          </tr>
                        </thead>
                        <tbody>
                          {row.original.child.map(child => (
                            <tr key={child.pid}>
                              <td>{child.pid}</td>
                              <td>{child.name}</td>
                              <td>{child.state}</td>
                            </tr>
                          ))}
                        </tbody>
                      </Table>
                    </td>
                  </tr>
                )}
              </React.Fragment>
            );
          })}
        </tbody>
      </Table>
    </div>
  );
}

export default TablaProcesos;
