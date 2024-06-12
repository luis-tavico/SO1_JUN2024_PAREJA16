import React, { useState, useEffect } from 'react';
import '../styles/Estilo.css'; // Asegúrate de crear este archivo para los estilos
import { FaSearch } from 'react-icons/fa'; // Importa el icono de lupa

function TablaProcesos() {
  const [processes, setProcesses] = useState([]);
  const [info, setInfo] = useState({});
  const [pid, setPid] = useState('');
  const [expandedRows, setExpandedRows] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    // Aquí llamarías a la API de Go para obtener los procesos.
    fetch('http://192.168.122.195:8080/procesos')
      .then(response => response.json())
      .then(data => {
        setProcesses(data.procesos);
        setInfo(data.info);
      })
      .catch(error => console.error('Error fetching processes:', error));
  }, []);

  const handleCreateProcess = () => {
    // Lógica para crear un proceso
    fetch('/procesos/crear', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
    })
    .then(response => response.json())
    .then(data => {
      setProcesses(data.procesos);
      setInfo(data.info);
    })
    .catch(error => console.error('Error creating process:', error));
  };

  const handleKillProcess = () => {
    // Lógica para matar un proceso
    fetch(`/procesos/eliminar/${pid}`, {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
    })
    .then(response => response.json())
    .then(data => {
      setProcesses(data.procesos);
      setInfo(data.info);
    })
    .catch(error => console.error('Error killing process:', error));
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

  const handleSearch = () => {
    // Lógica para buscar proceso
    // Aquí puedes usar el valor de searchTerm para buscar en la lista de procesos
    console.log('Searching for PID:', searchTerm);
  };

  const processesWithChildren = processes.filter(process => process.child && process.child.length > 0);
  const processesWithoutChildren = processes.filter(process => !process.child || process.child.length === 0);

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
            onChange={(e) => setSearchTerm(e.target.value)}
            placeholder="Buscar por PID"
          />
          <button onClick={handleSearch}><FaSearch /></button> {/* Botón de búsqueda con icono de lupa */}
        </div>
        <button onClick={handleKillProcess}>Matar Proceso</button>
      </div>
      <div className="table-container">
        {processesWithoutChildren.length > 0 && (
          <table>
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
              {processesWithoutChildren.map((process, index) => (
                <tr key={index}>
                  <td>{process.pid}</td>
                  <td>{process.nombre}</td>
                  <td>{process.estado}</td>
                  <td>{process.ram}</td>
                  <td>{process.usuario}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        {processesWithChildren.map((process, index) => (
          <div key={index} className="process-with-children">
            <table>
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
                <tr
                  onClick={() => toggleRow(index)}
                  className={expandedRows.includes(index) ? 'parent-row' : ''}
                >
                  <td>{process.pid}</td>
                  <td>{process.nombre}</td>
                  <td>{process.estado}</td>
                  <td>{process.ram}</td>
                  <td>{process.usuario}</td>
                </tr>
                {expandedRows.includes(index) && process.child && (
                  process.child.map((child, idx) => (
                    <tr key={`${index}-${idx}`} className="child-row">
                      <td>{child.pid}</td>
                      <td>{child.nombre}</td>
                      <td>{child.estado}</td>
                      <td>{child.ram}</td>
                      <td>{child.usuario}</td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        ))}
      </div>
    </div>
  );
}

export default TablaProcesos;


