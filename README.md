# Snippets

Find small snippet of codes written in golang and it can be deployed in AWS lambda to serve/solve a simple problem

## dailydigest

This snippet will trigger a API call at the given time across different timezone.

The body of the API call contains the following info: 
`
  {
    "timezone": "IST"
    "type" : "daily"
  }
`
