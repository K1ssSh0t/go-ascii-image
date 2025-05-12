package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg" // Registra el decodificador JPEG
	_ "image/png"  // Registra el decodificador PNG
	"log"
	"net/http" // Added for URL loading
	"os"
	"strings"

	"github.com/nfnt/resize" // Necesitarás instalar esta biblioteca: go get github.com/nfnt/resize
)

// Caracteres ASCII ordenados de oscuro a claro. Puedes ajustarlos.
//const asciiChars = "@%#*+=-:. " // Más denso

//const asciiChars = " `.-':_,^=;><+!rc*/z?sLTv)J7(|Fi{C}fI31tlu[neoZ5Yxjya]2ESwqkP6h9d4VpOGbUAKXHm8RD#$Bg0MNWQ%&@" // Más variado

// O una versión más simple y común:
//const asciiChars = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. "

// O incluso más simple:
const asciiChars = " .:!/r(l1Z4H9W8$@"
const resetColor = "\033[0m"
const blockChar = "█" // Character to use for color mode

func main() {
	// Definir flags para la línea de comandos
	imagePath := flag.String("image", "", "Ruta a la imagen de entrada (local o URL) (obligatorio)")
	outputFile := flag.String("output", "", "Ruta al archivo de salida para el arte ASCII (opcional)")
	maxWidth := flag.Uint("width", 80, "Ancho máximo del arte ASCII en caracteres") // Ancho común de terminal
	colorOutput := flag.Bool("color", false, "Habilitar salida a color en la consola")
	// También podrías añadir un flag para la altura o para invertir los caracteres

	flag.Parse()

	if *imagePath == "" {
		fmt.Println("Error: La ruta de la imagen es obligatoria.")
		flag.Usage() // Muestra cómo usar los flags
		os.Exit(1)
	}

	// Cargar la imagen
	var img image.Image
	var err error

	if strings.HasPrefix(*imagePath, "http://") || strings.HasPrefix(*imagePath, "https://") {
		img, err = loadImageFromURL(*imagePath)
	} else {
		img, err = loadImage(*imagePath)
	}

	if err != nil {
		log.Fatalf("Error al cargar la imagen: %v", err)
	}

	// Convertir a ASCII
	asciiArt := imageToASCII(img, *maxWidth, *colorOutput)

	// Imprimir en la consola
	fmt.Println(asciiArt)

	// Guardar en archivo si se especifica
	if *outputFile != "" {
		contentToSave := asciiArt
		if *colorOutput {
			// Si la salida a consola fue a color, generar una versión sin color para el archivo
			fmt.Println("\nGenerando versión sin color para el archivo...")
			contentToSave = imageToASCII(img, *maxWidth, false)
		}
		err := saveToFile(contentToSave, *outputFile)
		if err != nil {
			log.Fatalf("Error al guardar el archivo: %v", err)
		}
		fmt.Printf("\nArte ASCII guardado en: %s\n", *outputFile)
	}
}

// loadImageFromURL descarga y carga una imagen desde una URL.
func loadImageFromURL(url string) (image.Image, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("no se pudo descargar la imagen de '%s': %w", url, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error al descargar la imagen de '%s': status code %d", url, response.StatusCode)
	}

	img, format, err := image.Decode(response.Body)
	if err != nil {
		return nil, fmt.Errorf("no se pudo decodificar la imagen de la URL '%s': %w", url, err)
	}
	fmt.Printf("Imagen cargada desde URL: %s, Formato: %s\n", url, format)
	return img, nil
}

// loadImage carga una imagen desde la ruta especificada.
func loadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir el archivo '%s': %w", filePath, err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("no se pudo decodificar la imagen '%s': %w", filePath, err)
	}
	fmt.Printf("Imagen cargada: %s, Formato: %s\n", filePath, format)
	return img, nil
}

// imageToASCII convierte una imagen a una cadena de arte ASCII.
func imageToASCII(img image.Image, maxWidth uint, useColor bool) string {
	// Redimensionar la imagen para que se ajuste al ancho máximo, manteniendo la proporción
	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	// Calcular la nueva altura para mantener la relación de aspecto
	// Se multiplica por 0.5 (o un factor similar) porque los caracteres suelen ser más altos que anchos
	// y esto ayuda a que la proporción visual sea más correcta en la terminal.
	aspectRatio := float64(originalHeight) / float64(originalWidth)
	newHeight := uint(float64(maxWidth) * aspectRatio * 0.5) // Ajusta 0.4 según veas necesario

	// Usar la biblioteca de redimensionamiento
	resizedImg := resize.Resize(maxWidth, newHeight, img, resize.Lanczos3)
	resizedBounds := resizedImg.Bounds()
	width := resizedBounds.Dx()
	height := resizedBounds.Dy()

	var asciiBuilder strings.Builder

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Obtener el color del píxel
			pixel := resizedImg.At(x, y)
			r, g, b, _ := pixel.RGBA() // Obtiene RGBA, cada componente en [0, 65535]

			if useColor {
				// Convertir de [0, 65535] a [0, 255]
				r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
				// Usar ANSI escape code para el color del foreground
				// \033[38;2;r;g;bm establece el color RGB del primer plano
				asciiBuilder.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm", r8, g8, b8))
				asciiBuilder.WriteString(blockChar) // Puedes usar un carácter de bloque o el carácter ASCII
			} else {
				grayColor := color.GrayModel.Convert(pixel).(color.Gray)
				intensity := grayColor.Y // Intensidad del gris (0-255)

				// Mapear la intensidad al carácter ASCII
				charIndex := int(float64(intensity) / 255.0 * float64(len(asciiChars)-1))
				asciiBuilder.WriteByte(asciiChars[charIndex])
			}
		}
		if useColor {
			asciiBuilder.WriteString(resetColor) // Resetear color al final de la línea
		}
		asciiBuilder.WriteByte('\n') // Nueva línea al final de cada fila de píxeles
	}
	if useColor {
		asciiBuilder.WriteString(resetColor) // Asegurarse de que el color se resetea al final
	}

	return asciiBuilder.String()
}

// saveToFile guarda una cadena en un archivo.
func saveToFile(content string, filePath string) error {
	err := os.WriteFile(filePath, []byte(content), 0644) // 0644 son permisos estándar de archivo
	if err != nil {
		return fmt.Errorf("no se pudo escribir en el archivo '%s': %w", filePath, err)
	}
	return nil
}