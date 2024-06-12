import React, { useState, useEffect } from 'react';
import { useTable, useExpanded } from 'react-table';
import { Table } from 'react-bootstrap';
import { FaSearch } from 'react-icons/fa';
import '../components/Estilo.css';

function TablaProcesos() {
  const [processes, setProcesses] = useState([]);
  const [info, setInfo] = useState({});
  const [pid, setPid] = useState('');
  const [expandedRows, setExpandedRows] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    fetch('http://192.168.122.195:8080/procesos')
      .then(response => response.json())
      .then(data => {
        console.log('Data fetched:', data);  // Debugging output
        if (data && data.procesos) {
          setProcesses(data.procesos);
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
  }, []);

  const handleCreateProcess = () => {
    fetch('/procesos/crear', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
    })
    .then(response => response.json())
    .then(data => {
      if (data && data.procesos) {
        setProcesses(data.procesos);
        setInfo(data.info || {});
      }
    })
    .catch(error => console.error('Error creating process:', error));
  };

  const handleKillProcess = () => {
    fetch(`/procesos/eliminar/${pid}`, {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
    })
    .then(response => response.json())
    .then(data => {
      if (data && data.procesos) {
        setProcesses(data.procesos);
        setInfo(data.info || {});
      }
    })
    .catch(error => console.error('Error killing process:', error));
  };

  const handleSearchChange = event => {
    setSearchTerm(event.target.value);
  };

  const handleSearch = () => {
    console.log('Searching for PID:', searchTerm);
  };

  const toggleRow = (index) => {
    const currentIndex = expandedRows.indexOf(index);
    const newExpandedRows = [...expandedRows];
    if (currentIndex === -1) {
      newExpandedRows.push(index);
    } else {
      newExpandedRows.splice(currentIndex, 1);
    }
    setExpandedRows(newExpandedRows);
  };

  const columns = React.useMemo(
    () => [
      {
        Header: 'PID',
        accessor: 'pid',
      },
      {
        Header: 'Nombre',
        accessor: 'name',
      },
      {
        Header: 'Usuario',
        accessor: 'user',
      },
      {
        Header: 'Estado',
        accessor: 'state',
      },
      {
        Header: '%RAM',
        accessor: 'ram',
      },
      {
        Header: 'Acciones',
        id: 'expander',
        Cell: ({ row }) => (
          row.original.child && row.original.child.length > 0 && (
            <span {...row.getToggleRowExpandedProps()}>
              {row.isExpanded ? '▼' : '▶'}
            </span>
          )
        ),
      },
    ],
    []
  );

  const data = React.useMemo(
    () =>
      processes.map(process => ({
        ...process,
        subRows: process.child || [],
      })),
    [processes]
  );

  console.log('Columns:', columns);  // Debugging output
  console.log('Data:', data);        // Debugging output

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
    useExpanded
  );

  const filteredRows = rows.filter(row =>
    row.original.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="tabla-procesos-container">
      <div className="title-container">
        <h1>Tabla de Procesos</h1>
      </div>
      <div className="info-container">
        <p>En ejecución: {info.en_ejecucion}</p>
        <p>Suspendidos: {info.suspendidos}</p>
        <p>Detenidos: {info.detenidos}</p>
        <p>Zombies: {info.zombies}</p>
        <p>Total: {info.total}</p>
      </div>
      <div className="form-container">
        <button onClick={handleCreateProcess}>Crear Proceso</button>
        <div className="search-container">
          <input
            type="text"
            value={searchTerm}
            onChange={handleSearchChange}
            placeholder="Buscar por nombre"
          />
          <button onClick={handleSearch}><FaSearch /></button>
        </div>
        <input
          type="text"
          value={pid}
          onChange={(e) => setPid(e.target.value)}
          placeholder="PID para matar"
        />
        <button onClick={handleKillProcess}>Matar Proceso</button>
      </div>
      <div className="table-container">
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
            {filteredRows.map((row, index) => {
              prepareRow(row);
              return (
                <React.Fragment key={row.id}>
                  <tr {...row.getRowProps()} onClick={() => toggleRow(index)}>
                    {row.cells.map(cell => (
                      <td {...cell.getCellProps()}>{cell.render('Cell')}</td>
                    ))}
                  </tr>
                  {row.isExpanded && row.original.child && row.original.child.length > 0 && (
                    <tr>
                      <td colSpan={columns.length}>
                        <Table striped bordered hover>
                          <thead>
                            <tr>
                              <th>PID</th>
                              <th>Nombre</th>
                              <th>Estado</th>
                              <th>%RAM</th>
                              <th>Usuario</th>
                            </tr>
                          </thead>
                          <tbody>
                            {row.original.child.map(child => (
                              <tr key={child.pid}>
                                <td>{child.pid}</td>
                                <td>{child.name}</td>
                                <td>{child.state}</td>
                                <td>{child.ram}</td>
                                <td>{child.user}</td>
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
    </div>
  );
}

export default TablaProcesos;



