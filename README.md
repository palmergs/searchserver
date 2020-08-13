# searchserver

HTTP server to search for strings in documents. What is this good for? Not a whole lot, but if you have a fairly static set of tokens and a large number of document that you want to search to see if they contain those tokens this might be tool for you. Unfortunately, this code base is now about 5 years old and I wouldn't trust it any farther than I could throw it so use at your own risk.

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
 
  
