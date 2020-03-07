# What is this app?

This app routes all the requests that comes from the input to its respective game session. 
For example, the game is running on a Chrome tab and want to receive commands from a smartphone, actings as a controller.

#### How the communication occurs?

First, the game requests a new session, so players can join. The backend sends a code that identifies that session.
Second, the input controller sends the player information with the session code to the backend. This way, the backend service knows which player is playing wich game.
After that, all input controller actions is sent to the correct game session.

#### What communication technology is used?

All the communication between backend, game and input controller uses web sockets.

#### What is the public folder used for?

The public folder has just a application that simulates requests from the game and from the input controller. It's just to understande how the flow of requests works.


#### How to run the application?
##### Pre-requisites:
  - docker-compose

##### To run:
 
```
docker-compose up
```

Two ports will open on localhost: 8080 (backend service) and 3000 (input/game simulation service)


###### All game requests should be sent to /front namespace.
###### All input controller requests should be sent to /input namespace.


## Events from/to the game (/front):

- ###### create_session (game => backend)
```
"please, create session?"
```

- ###### enter_session (game => backend)
```
"1A2B"
```

- ###### session_entered (backend => game)
```
{ 
   "socketId":"3",
   "sessionCode":"QT6R",
   "players":[]
}
```

- ###### session_created (backend => game)
```
{ 
   "socketId":"3",
   "sessionCode":"QT6R",
   "players":[]
}
```

- ###### input_context (backend => game)
```
{ 
   "sessionCode":"PO6M",
   ...
}
```
*The backend will only need the sessionCode information to retrieve the front session and redirect this payload to it. Therefore, the payload can have **ANY** information in any form. The only constraint is that 1) the payload should be an object and 2) have a field called "sessionCode" of type string.*

- ###### session_changed (backend => game)
```
{ 
   "socketId":"3",
   "sessionCode":"QT6R",
   "players":[]
}
``` 

### Events from/to the controller (/input):

** REMEMBER: the backend just **

- ###### register_player (input => backend)
```
{ 
   "sessionCode":"ABCD",
   "player": {
      "id": "a2c79158-ea05-4a3d-adc0-27948b480e54",
      "name": "Zak",
      "color": "blue",
      "height": 1.9
   }
}
```

- ###### player_registered (backend => input)
```
"080ae447-2fa5-4e55-8d73-0c4669afc2d9"
```

- ###### input_context (input => backend)
```
{ 
   "sessionCode":"PO6M",
   "playerId":"a16dc6d7-a5cd-4775-82c2-1807ee6d9846",
   "activeActions":[ 
      "attack"
   ],
   "direction":{ 
      "x":0.1,
      "y":0.1
   }
   ...
}
```

*This input context is just an example of a valid payload. It can have any information in any form **AS LONG AS** it is an object and it has a "sessionCode" field of type string*