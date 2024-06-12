import React, { useState, useEffect } from 'react';
import { useTable, useExpanded } from 'react-table';
import { Table, Button, InputGroup, FormControl } from 'react-bootstrap';
import '../components/Estilo.css';

function TablaProcesos() {
  const [processes, setProcesses] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [searchPidTerm, setSearchPidTerm] = useState('');
  const [info, setInfo] = useState({});
  const [pidToKill, setPidToKill] = useState('');
  const [searchBy, setSearchBy] = useState('name'); 

  useEffect(() => {
    fetch('http://192.168.122.195:8080/procesos')
      .then(response => response.json())
      .then(data => {
        if (data && data.processes) {
          setProcesses(data.processes);
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

  const columns = React.useMemo(
    () => [
      {
        Header: 'PID',
        accessor: 'pid',
      },
      {
        Header: 'Name',
        accessor: 'name',
      },
      {
        Header: 'User',
        accessor: 'user',
      },
      {
        Header: 'State',
        accessor: 'state',
      },
      {
        Header: 'RAM',
        accessor: 'ram',
      },
      {
        Header: 'Actions',
        id: 'expander',
        Cell: ({ row }) => (
          <span {...row.getToggleRowExpandedProps()}>
            {row.isExpanded ? '▼' : '▶'}
          </span>
        ),
      },
    ],
    []
  );

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
    fetch(`/procesos/eliminar/${pidToKill}`, {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
    })
      .then(response => response.json())
      .then(data => {
        if (data && data.procesos) {
          setProcesses(data.procesos);
          setInfo(data.info || {});
        } else {
          alert('PID no encontrado');
        }
      })
      .catch(error => console.error('Error killing process:', error));
  };

  const data = React.useMemo(
    () =>
      processes.map(process => ({
        ...process,
        subRows: process.child || [],
      })),
    [processes]
  );

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

  const handleSearchChange = event => {
    setSearchTerm(event.target.value);
  };

  const handleSearchPidChange = event => {
    setSearchPidTerm(event.target.value);
  };

  const handleSearchByChange = () => {
    setSearchBy(prevSearchBy => (prevSearchBy === 'name' ? 'pid' : 'name'));
  };

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
