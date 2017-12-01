Llistar ordinadors en marxa de les aules
==============================================

Proves per generar un servei REST que  permeti fer servir JWT (JSON Web tokens) des de GO (golang).

> L'he intentat dividir en fitxers per fer-lo més semblant a com ho faria en Java (no m'acaba d'agradar el resultat)

Generació del programa
------------------------
Dependències:

* Faig servir diferents parts del framework web [Gorila](http://www.gorillatoolkit.org/): El mux, el context i els handlers (semblants als interceptors de Java)
* la implementació de JWT: [jwt-go](https://github.com/dgrijalva/jwt-go),  
* Una llibreria per convertir mapes a structs anomenada *mapstructure* (en realitat només em fa el codi més senzill):
* La configuració està en TOML i per tant fa falta la llibreria (go get github.com/naoina/toml)
* L'escanneig de la xarxa es fa amb la llibreria `listIP`(go get github.com/utrescu/listIP)

Per tant he configurat el projecte perquè faci servir 'dep' per instal·lar totes les dependències de forma automàtica. En el futur dep s'integrarà al paquet bàsic de Go, però per ara, cal instal·lar-lo:

    go get -u github.com/golang/dep/cmd/dep

Les dependències es descarregaran automàticament: 

    dep ensure

### Iniciar o compilar el programa

Un cop es tenen les dependències ja es pot iniciar el programa:

    go run *.go

També es pot aconseguir un binari nadiu del sistema en que s'executi el programa (en Linux l'executable generat agafa el nom del primer fitxer que troba en la llista)

    go build *.go
    ./aula

En Go també es poden generar executables de qualsevol plataforma. Per exemple podem generar un executable de Windows des de Linux:

    GOOS=windows GOARCH=amd64 go build *.go

(pot ser que tardi una mica més perquè ha d'aconseguir els binaris Windows)

### Accés al servei

Es pot anar amb el navegador a [http://localhost:3000](http://localhost:3000) per veure el servei en marxa.

Fa servir un frontend Angular per poder treballar amb el servei

![Angular](README/angular.png)

### Fitxer de configuració

#### Configuració bàsica

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

Aquí detecta els sistemes amb el port 22 obert (que és el que jo necessitava per control·lar les classes) però es pot posar un port diferent per cada classe o altres tipus de rangs

#### Base de dades

També fa servir una base de dades SQLite en la que s'hi defineix l'usuari i la contrasenya. El fitxer està en el directori de la configuració

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

He preparat una estructura de directoris per desplegar la part web */views/* per l'HTML.

L'autenticació JWT es fa enviant una capsalera 'Authorization'

> WARNING: El sistema no és segur si no es fa servir una connexió HTTPS.

Exemple d'ús amb un Authorization Token
--------------------------------------------

En les proves faré servir **httpie**

    pip install httpie

### Obtenir el token

Abans de poder fer servir els altres mètodes cal obtenir el token (en aquest moment funciona amb usuari "dept" i contrasenya "ies2017!"). En aquest cas envio 'pere'

    $ echo '{"username":"dept", "password":"ies2017!"}' | http http://localhost:3000/login

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
        "password": "ies2017!",
        "username": "dept"
    }

I ens donarà la resposta següent:

    HTTP/1.1 200 OK
    Content-Length: 173
    Content-Type: application/json
    Date: Wed, 01 Nov 2017 19:01:38 GMT

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

    { "aules": ["309","310","314"] }

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

També podem demanar per quins són els PC en marxa d'una classe. Per exemple la 309 podria tenir un resultat com aquest:

    $ http http://localhost:3000/aula/309/status Authorization:'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InBlcmUiLCJleHAiOjE1MDk1NzM3MTQsImlzcyI6ImxvY2FsaG9zdDozMDAwIn0.9LoBUzj4NaTHH8J02aWkR4DivJEvQFA2Pq8sHXMUtCk'

Que donarà:

    HTTP/1.1 200 OK
    Content-Length: 89
    Content-Type: application/json
    Date: Wed, 01 Nov 2017 21:42:29 GMT

    {
        "Aula": "309",
        "EnMarxa": [
            "192.168.9.22",
            "192.168.9.25",
            "192.168.9.26"
        ]
    }

TODO
================================

* Proporcionar alguna forma de crear usuaris
* Millorar el procediment d'escanneig dels ordinadors