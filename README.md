# searchserver

HTTP server to search for strings in documents

Usage:
    
    ./searchserver -f tokens.json
    
Query the `/tokens` endpoint to return a JSON document of all the loaded tokens.

Query the `/search` endpoint with a param `q` or a text body and return a JSON document of the location of all tokens found.

Example:

    GET http://localhost:6060/search?q=castle%20and%20road%20and%20house

Returns:

    [
      {
        "token": {
          "id": 3,
          "label": "road",
          "category": "noun"
        },
        "start_at": 11,
        "end_at": 15
      },
      {
        "token": {
          "id": 4,
          "label": "house",
          "category": "noun"
        },
        "start_at": 20,
        "end_at": 25
      }
    ]  
 
  
