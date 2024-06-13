import React, { useState, useEffect } from 'react'; // Importa React y hooks de estado y efectos
import { useTable, useExpanded } from 'react-table'; // Importa funciones de react-table para crear tablas y manejar filas expandidas
import { Table, Button, InputGroup, FormControl } from 'react-bootstrap'; // Importa componentes de react-bootstrap para la interfaz
import '../styles/EstiloProcesos.css'; // Importa estilos personalizados

function TablaProcesos() {
  // Define los estados del componente
  const [processes, setProcesses] = useState([]); // Estado para almacenar los procesos
  const [searchTerm, setSearchTerm] = useState(''); // Estado para el término de búsqueda por nombre
  const [searchPidTerm, setSearchPidTerm] = useState(''); // Estado para el término de búsqueda por PID
  const [info, setInfo] = useState({}); // Estado para la información adicional de los procesos
  const [pidToKill, setPidToKill] = useState(''); // Estado para almacenar el PID a eliminar
  const [searchBy, setSearchBy] = useState('name'); // Estado para determinar si se busca por nombre o por PID
  const [lastCreatedPid, setLastCreatedPid] = useState(''); // Estado para almacenar el PID del último proceso creado

  // Efecto para cargar los datos de los procesos desde la API
  useEffect(() => {
    // Función para cargar los datos de los procesos
    const fetchProcesses = () => {
      fetch('http://192.168.122.195:8080/procesos')
        .then(response => response.json())
        .then(data => {
          if (data && data.procesos.processes) {
            setProcesses(data.procesos.processes);
            setInfo(data.info || {});
          } else {
            setProcesses([]);
            setInfo({});
          }
        })
        .catch(error => {
          console.error('Error fetching processes:', error);
          setProcesses([]);
          setInfo({});
        });
    };

    // Llamar a fetchProcesses inicialmente y luego cada 2 segundos
    fetchProcesses(); // Llamar inicialmente

    const intervalId = setInterval(fetchProcesses, 2000); // Llamar cada 2 segundos

    // Limpiar intervalo al desmontar el componente
    return () => clearInterval(intervalId);
  }, []); // [] para ejecutar solo una vez al montar

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

  const handleCreateProcess = () => {
    fetch('http://192.168.122.195:8080/procesos/crear', {
      method: 'GET', // Cambia el método a GET según la definición de tu API
    })
      .then(response => response.json())
      .then(data => {
        if (data.estado === 'creado') {
          // Actualiza el estado con el PID del proceso creado
          setLastCreatedPid(data.pid);
          // Mostrar alerta flotante para proceso creado exitosamente
          showAlert(`Proceso creado con PID: ${data.pid}`);
        } else {
          // Mostrar alerta flotante para error al crear proceso
          showAlert('Error al crear proceso');
        }
      })
      .catch(error => {
        console.error('Error creating process:', error);
        // Mostrar alerta flotante para error general
        showAlert('Error al crear proceso');
      });
  };

  // Función para mostrar alertas flotantes
  const showAlert = (message) => {
    const alertContainer = document.createElement('div');
    alertContainer.className = 'alert-container';
    alertContainer.innerText = message;
    document.body.appendChild(alertContainer);

    setTimeout(() => {
      alertContainer.remove();
    }, 3000); // Desaparece después de 3 segundos (ajustable según necesidad)
  };

  const handleKillProcess = () => {
    fetch(`http://192.168.122.195:8080/procesos/eliminar/${pidToKill}`, {
      method: 'POST', // Asegúrate de usar el método correcto según tu API
      headers: { 'Content-Type': 'application/json' },
    })
      .then(response => response.json())
      .then(data => {
        if (data.estado === 'Eliminado') {
          // Mostrar alerta flotante para proceso eliminado exitosamente
          showAlert(`El proceso ${pidToKill} fue eliminado.`);
          // Verificar si el proceso eliminado es el último creado y actualizar el estado
          if (pidToKill === lastCreatedPid) {
            setLastCreatedPid('');
          }
        } else {
          // Mostrar alerta flotante para error al eliminar proceso
          showAlert('Error al eliminar proceso');
        }
      })
      .catch(error => {
        console.error('Error killing process:', error);
        // Mostrar alerta flotante para error general
        showAlert('Error al eliminar proceso');
      });
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
      <table className="horizontal-table">
        <tbody>
          <tr className="table-header">
            <td>En ejecución (1)</td>
            <td>Suspendidos (0)</td>
            <td>Detenidos (128)</td>
            <td>Zombies (1026)</td>
            <td>Total</td>
          </tr>
          <tr className='subRow'>
            <td>{info.en_ejecucion}</td>
            <td>{info.suspendidos}</td>
            <td>{info.detenidos}</td>
            <td>{info.zombies}</td>
            <td>{info.total}</td>
          </tr>
        </tbody>
      </table>

      <div>
        <Table className="no-border-table">
          <tbody>
          <tr>
            <td className="centered-cell">{lastCreatedPid ? `Último proceso creado: ${lastCreatedPid}` : 'Último proceso creado'}</td>
            <td>
              <Button variant="success" onClick={handleCreateProcess}>
                <i className="bi bi-file-earmark-plus"></i> Crear Proceso
              </Button>
            </td>
          </tr>
            <tr>
              <td>
                <FormControl
                  placeholder="PID para eliminar"
                  value={pidToKill}
                  onChange={(e) => setPidToKill(e.target.value)}
                />
              </td>
              <td>
                <Button variant="danger" onClick={handleKillProcess}>
                  <i className="bi bi-file-earmark-x"></i> Eliminar Proceso
                </Button>
              </td>
            </tr>
            <tr>
              <td>
                <FormControl
                  type="text"
                  placeholder={searchBy === 'name' ? 'Buscar por nombre' : 'Buscar por PID'}
                  value={searchBy === 'name' ? searchTerm : searchPidTerm}
                  onChange={searchBy === 'name' ? handleSearchChange : handleSearchPidChange}
                />
              </td>
              <td>
                <Button
                  style={{ fontSize: '14px', padding: '8px 12px', minHeight: '38px' }}
                  onClick={handleSearchByChange}
                >
                  <i className="bi bi-search"></i> {searchBy === 'name' ? 'Buscar por PID' : 'Buscar por Nombre'}
                </Button>
              </td>
            </tr>
          </tbody>
        </Table>
      </div>

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
