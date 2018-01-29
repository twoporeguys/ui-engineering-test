# Two Pore Guys UI Engineer test

## Goal
This exercise is meant to demonstrate candidate's ability to build a sample application, using a modern JS framework (ReactJS is preferred).
The design is not the main concern of the exercise, still the application should be good looking enough,  so that user's eyes doesn't start bleeding while using the application.

## Objective
The application will make use of a websocket endpoint to list wikipedia projects, list the pages in a given project and show changes on a page in realtime.
The websocket is `wss://wiki-meta-explorer.herokuapp.com/`.
It must be possible to:
- list the projects
- choose one project to see the list of its pages
- subscribe / unsubscribe to project changes (pages edition only) and to reflect those changes
- view the metadata about one specific page
- subscribe / unsubscribe to page changes (only one page can be subscribed to at a time)

It must be possible to browse through the projects / pages.

It is also recommended that you write at least one test as an exemplar of your approach to testing.

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
  "name": <NAME_OF_THE_COMMAND>,
  "data": <DATA_RETURNED_BY_THE_CALL>
}
```

Update messages have the form
```
{
  "id": <ID_MATCHING_THE_REQUEST_ID>,
  "name": <"page.update"|"project.update">,
  "data": <UPDATED_DATA> // ONLY CONTAINS CHANGED DATA
}
```

### Available commands

#### `project.list`
##### Description
Returns the list of projects available.
##### Arguments
Doesn't take any argument.

#### `page.list`
##### Description
Returns the list of pages in a given project.
##### Arguments
`project`: The name of the project as returned by `project.list` call.

#### `page.query`
##### Description
Return a specific page.
##### Arguments
`pageId`: The id of the page.

#### `project.subscribe`
##### Description
Subscribe to project updates.
##### Arguments
`project`: The name of the project as returned by `project.list` call.

#### `project.unsubscribe`
##### Description
Unsubscribe from project updates.
##### Arguments
Doesn't take any argument (there can be only one subscription a the time).

#### `page.subscribe`
##### Description
Subscribe to page updates.
##### Arguments
`pageId`: The id of the page.

#### `page.unsubscribe`
##### Description
Unsubscribe from page updates.
##### Arguments
Doesn't take any argument (there can be only one subscription a the time).

#### `ping`
##### Description
Ping the server. Useful to ensure the connection is kept open. The server might answer with `{name: "pong"}`.
##### Arguments
Doesn't take any argument.

### Data model
#### `page`
```
{
  "pageid": 9020,
  "title": "Daisy Duck",
  "pagelanguage": "en",
  "pagelanguagedir": "ltr",
  "length": 37954,
  "lastrevid": 797710133,
  "revisions": [
    {
      "revid": 797710133,
      "parentid": 797709597,
      "user": "XXX.YYY.ZZZ.WWW",
      "anon": "",
      "timestamp": "2017-08-28T19:03:35Z",
      "comment": ""
    }
  ]
}
```
