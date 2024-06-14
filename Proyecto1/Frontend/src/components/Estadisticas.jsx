import React, { useEffect, useState, useMemo } from 'react';
import { Doughnut } from 'react-chartjs-2';
import 'chart.js/auto'; // Importación automática de los componentes de Chart.js
import '../styles/Estilo.css'; // Asegúrate de crear este archivo para los estilos

function Estadisticas() {
  const [cpuData, setCpuData] = useState(null);
  const [ramData, setRamData] = useState(null);
  const url = "http://192.168.122.30:8080";

  useEffect(() => {
    const fetchData = () => {
      // Fetch data from the API
      fetch(url+'/estadisticas') // Reemplaza con tu endpoint real
        .then(response => response.json())
        .then(data => {
          setCpuData(parseFloat(data.cpu_percentage));
          setRamData(parseInt(data.ram_percentage));
          console.log('Datos recibidos:', data); // Imprimir en la consola
        })
        .catch(error => console.error('Error fetching data:', error));
    };

    fetchData(); // Realiza la primera llamada

    const interval = setInterval(() => {
      fetchData(); // Realiza una llamada cada 2 segundos
    }, 1000);

    return () => clearInterval(interval); // Limpia el intervalo al desmontar el componente
  }, []);

  const doughnutData = useMemo(() => (label, percentageUsed) => {
    const percentageFree = 100 - percentageUsed;
    return {
      labels: [`${label} Usado `, `${label} Libre`],
      datasets: [
        {
          data: [percentageUsed, percentageFree],
          backgroundColor: ['#FF6384', '#36A2EB'],
        },
      ],
    };
  }, []);

  return (
    <div className="estadisticas-container">
      <div className="title-container">
        <h1>SO1 - JUN 2024</h1>
      </div>
      <div className="charts-container">
        {cpuData !== null && (
          <div className="chart">
            <h2>% CPU</h2>
            <Doughnut data={doughnutData('CPU', cpuData)} />
          </div>
        )}
        {ramData !== null && (
          <div className="chart">
            <h2>% RAM</h2>
            <Doughnut data={doughnutData('RAM', ramData)} />
          </div>
        )}
      </div>
    </div>
  );
}

export default Estadisticas;
