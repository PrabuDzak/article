
dzakshort.com


shortenId -> alsknc. alphanumeric. exactly 6 char.

API
- generate url - dzakshort.com/short
  {
    url: "www.google.com"
  }
  resp
  {
    id: "asca"
    shortUrl: "dzakshort.com/s/alsknc"
  }

- redirect - dzakshort.com/s/{shortenId}
  - dzakshort.com/s/alsknc -> www.google.com

- stats - dzakshort.com/stats/{shortenId}
  {
    shortUrl: "dzakshort.com/s/alsknc"
  }
  {
    counter: 
    createdAt:
  }

idGenerator -> "alsknc"
- uuid -> truncate
- hash(time) -> truncate
