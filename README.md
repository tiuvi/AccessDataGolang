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

buffer de bytes
Dos escrituras bufferBytes en:           0. 000 044 401 segundos.
Dosmil escrituras bufferBytes en:        0. 168 038 051 segundos.
veintemil escrituras bufferBytes en:     0. 463 238 103 segundos.
doscientasmil escrituras bufferBytes en: 7. 530 359 778 segundos.

buffer de bits
Dosmil escrituras bufferBytes en:         0. 075 561 161
veintemil escrituras bufferBytes en:      1. 246 176 394
doscientasmil escrituras bufferBytes en:  7. 409 574 141


Canal de buffer de bytes
Dos escrituras ChanBufferBytes en:            0. 000 066 830
Dosmil escrituras ChanBufferBytes en:         0. 051 327 189
Veintemil escrituras ChanBufferBytes en:      0. 996 522 694
doscientasmil escrituras ChanBufferBytes en: 12. 643 067 397

Canal de buffer de bits
Dosmil escrituras ChanBufferBytes en:        	0. 114 776 100
Veintemil escrituras ChanBufferBytes en:        1. 365 854 612
doscientasmil escrituras ChanBufferBytes en:   12. 303 100 852


Las escrituras en mapas son en dos columnas pero lo escribe en una vez.
Una escrituras MapsBufferBytes en:       0. 000 081 983
Mil escrituras MapsBufferBytes en:       0. 078 245 810
Diezmil escrituras ChanBufferBytes en:   0. 931 418 163
Cienmil escrituras ChanBufferBytes en:  10 .808 210 071

mapa buffer de bits
Mil escrituras MapsBufferBytes en:       0. 173 969 542
Diezmil escrituras ChanBufferBytes en:   1. 930 629 185
Cienmil escrituras ChanBufferBytes en:  10. 658 817 437
