# ASCII Art Generator

Este programa convierte imágenes (locales o desde una URL) en arte ASCII.

## Requisitos

- Go (lenguaje de programación)
- La biblioteca `github.com/nfnt/resize`. Puedes instalarla con:
  ```bash
  go get github.com/nfnt/resize
  ```

## Compilación

Para compilar el programa, navega al directorio donde se encuentra `ascii-art.go` y ejecuta:

```bash
go build ascii-art.go
```

Esto generará un archivo ejecutable llamado `ascii-art` (o `ascii-art.exe` en Windows).

## Uso

Puedes ejecutar el programa desde la línea de comandos.

### Sintaxis

```bash
./ascii-art -image <ruta_o_url_imagen> [opciones]
```

### Flags Obligatorios

-   `-image <ruta_o_url_imagen>`: Especifica la ruta a un archivo de imagen local o una URL de una imagen.

### Flags Opcionales

-   `-output <ruta_archivo_salida>`: Guarda el arte ASCII generado en el archivo especificado. Si no se proporciona, el arte ASCII se imprimirá en la consola.
-   `-width <numero>`: Establece el ancho máximo del arte ASCII en caracteres. El valor predeterminado es `80`. La altura se ajustará automáticamente para mantener la proporción de la imagen.

### Ejemplos

1.  **Convertir una imagen local e imprimir en la consola:**

    ```bash
    ./ascii-art -image ./mi_imagen.jpg
    ```

2.  **Convertir una imagen desde una URL y guardarla en un archivo:**

    ```bash
    ./ascii-art -image https://www.ejemplo.com/imagen_online.png -output arte_ascii.txt
    ```

3.  **Convertir una imagen local con un ancho personalizado:**

    ```bash
    ./ascii-art -image ./otra_imagen.jpeg -width 120
    ```

4.  **Convertir una imagen desde una URL, con ancho personalizado y guardarla en un archivo:**
    ```bash
    ./ascii-art -image https://www.ejemplo.com/logo.png -width 50 -output logo_ascii.txt
    ```

## Funcionamiento Interno

1.  **Carga de Imagen**:
    *   Si la ruta proporcionada comienza con `http://` o `https://`, el programa descarga la imagen desde la URL.
    *   De lo contrario, intenta cargar la imagen desde el sistema de archivos local.
    *   Soporta formatos de imagen comunes como JPEG y PNG.
2.  **Redimensionamiento**:
    *   La imagen se redimensiona al ancho especificado (o al predeterminado de 80 caracteres).
    *   La altura se calcula para mantener la relación de aspecto original, con un ajuste para compensar que los caracteres de la terminal suelen ser más altos que anchos.
3.  **Conversión a Escala de Grises**:
    *   Cada píxel de la imagen redimensionada se convierte a su equivalente en escala de grises.
4.  **Mapeo a Caracteres ASCII**:
    *   La intensidad de cada píxel en escala de grises (0-255) se mapea a un carácter de la cadena `asciiChars` definida en el código. Los píxeles más oscuros se mapean a caracteres que visualmente parecen más densos, y los más claros a caracteres más dispersos.
5.  **Salida**:
    *   El resultado es una cadena de texto que representa la imagen en arte ASCII.
    *   Esta cadena se imprime en la consola y/o se guarda en un archivo si se especifica la opción `-output`.

## Personalización de Caracteres ASCII

Puedes modificar la constante `asciiChars` en el archivo `ascii-art.go` para cambiar el conjunto de caracteres utilizado para generar el arte. Hay varias opciones comentadas en el código, desde conjuntos más densos hasta más simples. Experimenta para encontrar el que mejor se adapte a tus preferencias. Recuerda que los caracteres deben estar ordenados aproximadamente de oscuro/denso a claro/disperso.