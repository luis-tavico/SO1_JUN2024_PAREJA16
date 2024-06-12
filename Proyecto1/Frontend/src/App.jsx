import { useState } from 'react'
import reactLogo from './assets/react.svg'
import bootstrapLogo from './assets/Bootstrap-logo.png'
import viteLogo from './assets/vite.svg'
import ModoOscuro from "./components/ModoOscuro"
import './styles/App.css'
import Head from "./components/Head"
import Estadisticas from './components/Estadisticas'
import TablaProcesos from './components/TablaProcesos'


function App() {
  let component
  switch (window.location.pathname) {
    case "/":
      component = <Estadisticas />
      // Lógica para la ruta "/estadisticas"
      break;
    case "/estadisticas":
      component = <Estadisticas />
      // Lógica para la ruta "/estadisticas"
      break;
    case "/tablaprocesos":
      component = <TablaProcesos />
      break;
    // Otros casos aquí
    default:
      // Lógica para otras rutas
  }

  return (
    <>
    <Head />
    {component}
    </>
  )
  
}

export default App
