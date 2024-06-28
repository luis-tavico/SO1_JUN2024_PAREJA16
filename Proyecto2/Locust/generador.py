import random
import json

# Lista extendida de frases relacionadas con el clima
frases_clima = [
    "Lindo día soleado",
    "Agradable clima fresco",
    "Me encanta la lluvia",
    "Hace mucho viento hoy",
    "Tormenta intensa esta noche",
    "Granizo inesperado",
    "Niebla densa por la mañana",
    "Calor sofocante",
    "Frío intenso esta semana",
    "Huracán en la región",
    "Brisa suave y agradable",
    "Clima cambiante constantemente",
    "Aroma a tierra mojada después de la lluvia",
    "Paisaje cubierto de nieve",
    "Viento frío del norte",
    "Tarde soleada y tranquila",
    "Frente frío acercándose",
    "Nubes oscuras de tormenta",
    "Nebulosa mañana de otoño",
    "Atardecer con colores intensos",
    "Ráfagas de viento fuerte",
    "Olas altas en la costa",
    "Noche estrellada y despejada",
    "Nieve fresca sobre las montañas",
    "Lluvia ligera y persistente",
    "Mañana húmeda y fresca",
    "Calor sofocante en la ciudad",
    "Heladas matinales en invierno",
    "Tormenta de verano repentina",
    "Niebla espesa en el valle"
]

# Lista más extensa de países
paises = [
    "Argentina", "Australia", "Brasil", "Canadá", "Chile", "China", "Colombia", "Egipto", "Francia", "Alemania", 
    "India", "Indonesia", "Italia", "Japón", "Corea del Sur", "México", "Países Bajos", "Nueva Zelanda", 
    "Perú", "Filipinas", "Rusia", "Arabia Saudita", "España", "Suecia", "Suiza", "Tailandia", "Turquía", 
    "Ucrania", "Reino Unido", "Estados Unidos", "Afganistán", "Albania", "Angola", "Barbados", "Belice", 
    "Benín", "Bolivia", "Burkina Faso", "Burundi", "Cabo Verde", "Camboya", "Chad", "Costa Rica", "Cuba", 
    "Dinamarca", "Ecuador", "Etiopía", "Fiyi", "Finlandia", "Ghana", "Grecia", "Honduras", "Irán", "Iraq", 
    "Jamaica", "Kenia", "Laos", "Letonia", "Liberia", "Madagascar", "Malasia", "Mali", "Mauritania", "Mongolia", 
    "Mozambique", "Namibia", "Nepal", "Nicaragua", "Noruega", "Omán", "Pakistán", "Panamá", "Papúa Nueva Guinea", 
    "Qatar", "Ruanda", "Senegal", "Siria", "Somalia", "Sudáfrica", "Sudán", "Surinam", "Tanzania", "Uganda", 
    "Vietnam", "Yemen", "Zambia", "Zimbabue"
]

# Función para generar datos aleatorios
def generar_datos(num_datos):
    datos = []
    for _ in range(num_datos):
        texto = random.choice(frases_clima)
        pais = random.choice(paises)
        datos.append({"texto": texto, "pais": pais})
    return datos

# Generar 10,000 datos aleatorios
datos_json = generar_datos(10000)

# Guardar los datos en un archivo JSON
with open('datos_clima_paises.json', 'w', encoding='utf-8') as file:
    json.dump(datos_json, file, ensure_ascii=False, indent=4)

print("Archivo JSON generado exitosamente.")
