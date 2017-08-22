# Two Pore Guys UI Engineer test

## Goal
This exercise is meant to demonstrate candidate's ability to build a sample application, using a modern JS framework (ReactJS is preferred).
The design is not the main concern of the exercise, still the application should be good looking enough,  so that user's eyes doesn't start bleeding while using the application.

## Objective
The application will make use of a websocket endpoint to list city bikes services, list the stations of one of those services with their status, available bikes and free slots.
The websocket is wss://find-a-bike.herokuapp.com/.
It must be possible to list the networks, to choose one to see the list of its stations, to subscribe / unsubscribe to stations changes and to reflect those changes.
The stations should be grouped by department and city.
Currently, he only available network in the backend is `'velib'`. Choosing another one should be either prevented and / or should display an explanatory message to the user.

## Backend API

### Messages format
Messages can have three types:
- Request: From client to server
- Response: From server to client, as a request answer
- Update: From server to client, notify that data have changed on server side

Request messages have the form
```
{
  "id": <ANY_UNIQUE_ID>,
  "name": <NAME_OF_THE_COMMAND>,
  "args": { // OPTIONAL PROPERTY
    <KEY>: <VALUE>
  }
}
```

Response messages have the form
```
{
  "id": <ID_MATCHING_THE_REQUEST_ID>,
  "name": "response",
  "data": <DATA_RETURNED_BY_THE_CALL>
}
```

Update messages have the form
```
{
  "id": <ID_MATCHING_THE_REQUEST_ID>,
  "name": "update",
  "type": <TYPE_OF_UPDATED_DATA>,
  "data": <UPDATED_DATA> // ONLY CONTAINS CHANGED DATA
}
```

### Available commands

#### `network.list`
##### Description
Returns the list of networks available.
##### Arguments
Doesn't take any argument.

#### `station.list`
##### Description
Returns the list of stations available for a given network.
##### Arguments
`network`: The id of the network. So far, the only acceptable value is `'velib'`. Any other value will never return.

#### `station.subscribe`
##### Description
Subscribe to stations updates.
##### Arguments
`network`: The id of the network. So far, the only acceptable value is `'velib'`. Any other value will never return nor send any update.

#### `station.unsubscribe`
##### Description
Subscribe to stations updates.
##### Arguments
Doesn't take any argument. The latest subscription is canceled.

### Data model
#### `network`
```
{
  "id": <NETWORK_ID>,
  "name": <DISPLAY_NAME_OF_THE_NETWORK>
  "country": <COUNTRY_WHERE_THE_NETWORK_OPERATE>,
  "city": <CITY_WHERE_THE_NETWORK_OPERATE>
}
```

#### `station`
```
{
  "id": <STATION_ID>,
  "name": <DISPLAY_NAME_OF_THE_STATION>
  "empty_slots": <NUMBER_OF_FREE_SLOTS_IN_THE_STATION>,
  "free_bikes": <NUMBER_OF_AVAILABLE_BIKES_IN_THE_STATION>,
  "department": <CODE_OF_THE_COUNTY_THE_STATION_IS_LOCATED_IN>
  "city": <CITY_THE_STATION_IS_LOCATED_IN>
}
```
