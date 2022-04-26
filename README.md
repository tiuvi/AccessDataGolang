# AccessDataGolang

Nueva tecnología de acceso a datos, separa la busqueda de datos, del acceso a datos aumentando la velocidad en *100.

Sin inyeccion sql, es totalmente imposible inyectar codificación de cualquier manera ademas tiene un prefiltrado y un 
postfiltrado en la entrada y la salida de datos del archivo.

Se pierde capacidad de memoria a base de ganar rendimiento de CPU, si tu cuello de botella es la memoria entonces el 
acceso a datos no es para ti.

El acceso a datos es una solución pensada para crear una red social con muy poco recursos que pueda competir con redes grandes
como Facebook, twitter, youtube, tiktok.

Solucion pensada para todos los formatos de contenido video,audio, imagenes.

Atomización de los datos, puedes tener en un solo archivo un video, los likes del video, un contador de visitas del video...etc

Provado con mas de 1000 campos estaticos + 1000 columnas dinamicas.

Si crecen los datos tu velocidad de acceso a ellos no decrece totalmente lineal, ademas compatible en concurrencia y en paralelo.
 <br> <br>
buffer de bytes <br>
Dos escrituras bufferBytes en:           0. 000 044 401 segundos. <br>
Dosmil escrituras bufferBytes en:        0. 168 038 051 segundos. <br>
veintemil escrituras bufferBytes en:     0. 463 238 103 segundos. <br>
doscientasmil escrituras bufferBytes en: 7. 530 359 778 segundos. <br>
<br>
buffer de bits <br>
Dosmil escrituras bufferBytes en:         0. 075 561 161 <br>
veintemil escrituras bufferBytes en:      1. 246 176 394 <br>
doscientasmil escrituras bufferBytes en:  7. 409 574 141 <br>

<br>
Canal de buffer de bytes <br>
Dos escrituras ChanBufferBytes en:            0. 000 066 830 <br>
Dosmil escrituras ChanBufferBytes en:         0. 051 327 189 <br>
Veintemil escrituras ChanBufferBytes en:      0. 996 522 694 <br>
doscientasmil escrituras ChanBufferBytes en: 12. 643 067 397 <br>
<br>
Canal de buffer de bits <br>
Dosmil escrituras ChanBufferBytes en:         	0. 114 776 100 <br>
Veintemil escrituras ChanBufferBytes en:        1. 365 854 612 <br>
doscientasmil escrituras ChanBufferBytes en:   12. 303 100 852 <br>


Las escrituras en mapas son en dos columnas pero lo escribe en una vez. <br>
Una escrituras MapsBufferBytes en:       0. 000 081 983 <br>
Mil escrituras MapsBufferBytes en:       0. 078 245 810 <br>
Diezmil escrituras ChanBufferBytes en:   0. 931 418 163 <br>
Cienmil escrituras ChanBufferBytes en:  10 .808 210 071 <br>
 <br>
mapa buffer de bits <br>
Mil escrituras MapsBufferBytes en:       0. 173 969 542 <br>
Diezmil escrituras ChanBufferBytes en:   1. 930 629 185 <br>
Cienmil escrituras ChanBufferBytes en:  10. 658 817 437 <br>
