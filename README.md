Llistar ordinadors en marxa de les aules
==============================================

Proves per generar un servei REST que  permeti fer servir JWT (JSON Web tokens) des de GO (golang).

> L'he intentat dividir en fitxers per fer-lo més semblant a com ho faria en Java (no m'acaba d'agradar el resultat)

Generació del programa
------------------------

Faig servir diferents parts del framework web [Gorila](http://www.gorillatoolkit.org/): El mux, el context i els handlers (semblants als interceptors de Java)

    go get github.com/gorilla/context
    go get github.com/gorilla/handlers
    go get github.com/gorilla/mux

La implementació de JWT: [jwt-go](https://github.com/dgrijalva/jwt-go)

    go get github.com/dgrijalva/jwt-go

I una llibreria per convertir mapes a structs anomenada *mapstructure* (en realitat només em fa el codi més senzill):

    go get github.com/mitchellh/mapstructure

La configuració està en TOML i per tant fa falta la llibreria

    go get github.com/naoina/toml

L'escanneig de la xarxa es fa amb la llibreria `listIP`:

    go get github.com/utrescu/listIP

### Iniciar o compilar el programa

Després només fa falta iniciar el programa:

    go run *.go

També es pot aconseguir un binari nadiu del sistema en que es compili (en Linux l'executable generat agafa el nom del primer fitxer que troba en la llista)

    go build *.go
    ./aula

En Go es poden generar executables de qualsevol plataforma. Per exemple podem generar un executable de Windows des de Linux:

    GOOS=windows GOARCH=amd64 go build *.go

(pot ser que tardi una mica més perquè ha d'aconseguir els binaris Windows)

### Accés al servei

Es pot anar amb el navegador a [http://localhost:3000](http://localhost:3000) per veure el servei en marxa. Inicialment s'accedeix a un formulari de login que serveix per iniciar-se en el sistema (però la idea no és fer una aplicació web sinó una API REST)

### Fitxer de configuració

El programa carrega les dades de les classes d'un fitxer en format TOML

Per tant cal tenir un fitxer amb la configuració de les aules en el mateix directori del binari.

Per exemple aquesta seria la configuració de dues aules 309 i 310:

    [aules]

      [aules.309]
      rang = "192.168.9.0/24"
      name = "Aula 309"
      port = 22

      [aules.310]
      rang = "192.168.10.0/24"
      name = "Aula 310"
      port = 22

Aquí detecta els sistemes amb el port 22 obert (que és el que jo necessitava per control·lar les classes) però es pot posar un port diferent per cada classe

Descripció del servei
------------------------

La base del programa és una interfície REST que retorna els resultats en format JSON.

Per la interfície REST he definit diferents rutes en el programa (però no n'hi ha cap que faci res d'interessant, només són l'esboç d'una idea del que vull implementar en el futur)

| URL                   | Mètode  |  Funció                                                  |
|-----------------------|---------|----------------------------------------------------------|
| /login                | POST    | S'hi envia l'usuari i la contrasenya en el body JSON. Retorna el token a fer servir |
| /aula/list            | GET     | Llista les aules (necessita el token o donarà error)     |
| /aula/{numero}/status | GET     | Llista les màquines en marxa (necessita el token)        |
| /aula/{numero}/stop   | POST    | Encara no fa res ... (necessita el token o donarà error) |

Un valor important i que s'hauria de mantenir en secret és la clau de xifrat que es fa servir per generar els tokens:

    var clauDeSignat = []byte("SiLaLletFosXocolataNoCaldriaColacao")

He preparat una estructura de directoris per desplegar la part web */views/* per l'HTML i */static/* pels recursos.

Faig servir dues formes d'autenticació:

* Enviant una capsalera 'Authorization'
* Fent servir una cookie

> WARNING: Cap dels sistemes és segur si no es fa servir una connexió HTTPS. En el cas de les Cookies s'han de restringir perquè no es puguin llegir des de Javascript

Exemple d'ús amb un Authorization Token
--------------------------------------------

En les proves faré servir **httpie**

    pip install httpie

### Obtenir el token

Abans de poder fer servir els altres mètodes cal obtenir el token (en aquest moment funciona amb qualsevol usuari i contrasenya). En aquest cas envio 'pere'

    $ echo '{"username":"pere", "password":"contra"}' | http http://localhost:3000/login

Amb aquesta comanda es generarà una petició POST:

    POST /login HTTP/1.1
    Accept: application/json
    Accept-Encoding: gzip, deflate
    Connection: keep-alive
    Content-Length: 38
    Content-Type: application/json
    Host: localhost:3000
    User-Agent: HTTPie/0.9.4

    {
        "password": "contra",
        "username": "pere"
    }

I ens donarà la resposta següent:

    HTTP/1.1 200 OK
    Content-Length: 173
    Content-Type: application/json
    Date: Wed, 01 Nov 2017 19:01:38 GMT
    Set-Cookie: Auth=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InBlcmUiLCJleHAiOjE1MDk1NjY0OTgsImlzcyI6ImxvY2FsaG9zdDozMDAwIn0.FlGUcQMG6U4c7yWIhS3QwDC5ervictvHfThGph7d4s4; Expires=Wed, 01 Nov 2017 20:01:38 GMT; HttpOnly

    {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InBlcmUiLCJleHAiOjE1MDk1NjY0OTgsImlzcyI6ImxvY2FsaG9zdDozMDAwIn0.FlGUcQMG6U4c7yWIhS3QwDC5ervictvHfThGph7d4s4"
    }

El valor de *token* és el que necessitem per mantenir l'autenticació (al cap d'una hora caduca). 

### Llistar les classes

Amb el token es poden fer peticions a les altres URL. Només cal posar-lo en la capsalera de la petició: **Authorization*:

    $ http http://localhost:3000/aula/list Authorization:'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InBlcmUiLCJleHAiOjE1MDk1NjY0OTgsImlzcyI6ImxvY2FsaG9zdDozMDAwIn0.FlGUcQMG6U4c7yWIhS3QwDC5ervictvHfThGph7d4s4'

Que donarà aquesta resposta:

    HTTP/1.1 200 OK
    Content-Length: 114
    Content-Type: application/json
    Date: Wed, 01 Nov 2017 19:47:30 GMT

    ["309","310","314"]

El Token és vàlid durant una hora. Per tant si repetim la petició després d'aquest temps el resultat serà que ja no podem identificar-nos:

    HTTP/1.1 200 OK
    Content-Length: 31
    Content-Type: application/json
    Date: Wed, 01 Nov 2017 20:08:42 GMT

    {
        "message": "Token is expired"
    }

En cas de que es faci la petició sense token també rebrem un error:

    HTTP/1.1 200 OK
    Content-Length: 50
    Content-Type: application/json
    Date: Wed, 01 Nov 2017 20:36:38 GMT

    {
        "message": "An authorization header is required"
    }

### Demanar pels PC en marxa d'una classe

També podem demanar per quins són els PC en marxa d'una classe. Per exemple la 309:

    $ http http://localhost:3000/aula/309/status Authorization:'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InBlcmUiLCJleHAiOjE1MDk1NzM3MTQsImlzcyI6ImxvY2FsaG9zdDozMDAwIn0.9LoBUzj4NaTHH8J02aWkR4DivJEvQFA2Pq8sHXMUtCk'

Que donarà:

    HTTP/1.1 200 OK
    Content-Length: 59
    Content-Type: application/json
    Date: Wed, 01 Nov 2017 21:42:29 GMT

    {
        "Aula": "309",
        "EnMarxa": [
            "i309-01m",
            "i309-01d",
            "i309-03e"
        ]
    }

Exemple d'ús amb Cookies
--------------------------------------------

Si es fan servir les Cookies es pot treballar des del navegador. Primer s'accedeix a la pantalla de login:

    http://localhost:3000

![login](README/login.png)

Que retornarà el token (innecessari perquè també està en una cookie)

![login token](README/login2.png)

A partir d'aquest moment es pot navegar per les adreces protegides sense problemes (el token està en la cookie)

![aules](README/aules.png)

![aula309](README/aula309.png)

Es pot tancar la sessió accedint a la URL [http://localhost:3000/logout](http://localhost:3000/logout)

I si tornem a intentar anar a pàgines protegides donarà error:

![Error](README/error.png)
