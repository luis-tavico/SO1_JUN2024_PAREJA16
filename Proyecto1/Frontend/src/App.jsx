import { useState } from 'react'
import reactLogo from './assets/react.svg'
import bootstrapLogo from './assets/Bootstrap-logo.png'
import viteLogo from './assets/vite.svg'
import ModoOscuro from "./components/ModoOscuro"
import './styles/App.css'
import Head from "./components/Head"
import Estadisticas from './components/Estadisticas'
import Procesos from './components/Procesos'


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
    case "/procesos":
      component = <Procesos />
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
