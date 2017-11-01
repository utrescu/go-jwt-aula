Proves per generar un servei REST que  permeti fer servir JWT (JSON Web tokens) des de GO (golang).

> L'he intentat dividir en fitxers per fer-lo més semblant a com ho faria en Java (no m'acaba d'agradar el resultat)

Faig servir diferents parts de [Gorila](http://www.gorillatoolkit.org/). El mux i els handlers (coneguts com interceptors en Java)

    $ go get github.com/gorilla/context
	$ go get github.com/gorilla/handlers
	$ go get github.com/gorilla/mux

La implementació de JWT: [jwt-go](https://github.com/dgrijalva/jwt-go)

    $ go get github.com/dgrijalva/jwt-go

I una llibreria per convertir mapes a structs anomenada *mapstructure* (en realitat només em fa el codi més senzill):

    $ go get github.com/mitchellh/mapstructure

Després només fa falta iniciar el programa:

    $ go run *.go

També es pot aconseguir un binari nadiu compilant-lo (en Linux l'executable generat agafa el primer que troba en la llista)

    $ go build *.go
    $ ./aula

Amb el navegador a http://localhost:3000 s'accedeix a la pàgina inicial

Implementació
-------------------

La base del programa és una interfície REST que retorna els resultats en format JSON.

Per la interfície REST he definit diferents rutes en el programa (però no n'hi ha cap que faci res d'interessant, només són l'esboç d'una idea del que vull implementar en el futur)

| URL                   | Mètode  |  Funció                                |
|-----------------------|---------|----------------------------------------|
| /login                | POST    | S'hi envia l'usuari i la contrasenya en el body JSON. Retorna el token a fer servir |
| /aula/list            | GET     | Llista les aules (necessita el token o donarà error)    |
| /aula/{numero}/status | GET     | Llista les màquines en marxa (necessita el token)       |
| /aula/{numero}/stop   | POST    | Encara no fa res ... (necessita el token o donarà error)  |

Un valor important i que s'hauria de mantenir en secret és la clau de xifrat que es fa servir per generar els tokens:

    var clauDeSignat = []byte("SiLaLletFosXocolataNoCaldriaColacao")

He preparat una estructura de directoris per desplegar la part web */views/* per l'HTML i */static/* pels recursos.

Exemple d'ús
----------------------------

En les proves faré servir **httpie**

    $ pip install httpie

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

    [
        {
            "Nom": "309",
            "Numero": 309,
            "Xarxa": "192.168.9.0/24"
        },
        {
            "Nom": "310",
            "Numero": 310,
            "Xarxa": "192.168.10.0/24"
        },
        {
            "Nom": "314",
            "Numero": 314,
            "Xarxa": "192.168.9.16/24"
        }
    ]

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

### Demanar pels PC d'una classe 

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