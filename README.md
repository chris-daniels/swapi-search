# Star Wars Character Search CLI

A quick integration with https://swapi.dev/api to search characters. 

### Requirements
 - Go version 1.22.0. 

### To run
```
$ go run .
```

### To test
```
$ go test .
```

### My approach
SWAPI supports a search term on each entity type - planets, films, people, etc. So I didn't see any way around making a search call for each of these types. 

I decided to cache each of these SWAPI search calls for a given search term in-memory. I used a [caching library](https://github.com/patrickmn/go-cache) with a TTL of one hour. 

After fetching each entity, we're provided an array of people URLs. I just go and fetch those individually, surrounded by another cache. 

This presents a problem, though. A broad enough search term will require us to go fetch every person in the system individually, which is very slow. To get around this, I make a series of paginated `/people` requests _at startup_ to get all people entries cached in memory. This takes a few seconds up front, but it improves the performance of search after that significantly. 

The project is pretty small, but it's roughly split as follows:
 - `main.go` handles startup and reading user input
 - `search.go` handles the logic of taking a search term, fanning that search term out to a bunch of requests, and building a text response to be surfaced to the user. 
 - `swapi.go` handles making requests to the SWAPI API.


 There are some basic tests in `main_test.go` that mock the SWAPI API. 


### Bringing this to production
To bring this app to production, we would need to make a few changes. 

We could distribute this CLI tool, and people could use it that way, but eventually we might want something prettier. Also, caching is done by each instance of the app here, which means the user will have to wait around for a few seconds at startup to populate the people cache. They'll also get pretty slow responses for every new query they do. For these reasons we might want to centralize some stuff.

So to move forward, I'd imagine we would have a client/server model in place here. The server would handle fetching data from the SWAPI and providing a search API, and clients would then handle taking user input and making search requests. This could still be a CLI tool, but we could also consider building a web front-end. 

If we want to scale our server to handle many, many search requests, we'll want to scale horizontally. In this case, we probably want to share the cache among instances of this service. We might swap out our in-memory solution here for a fast key-value store that's shared. Maybe we could use a Redis instance here. 

One other thing I might consider for this use case is a different approach, where we periodically pull the entire SWAPI dataset (it's quite small and can easily fit in memory), and build out our own in-memory search index. The logic might be a little tricky, but I think we could set up some sort of tree structure that allows us to fuzzy search on all the entities in the database. There are only a few hundred entities in the Star Wars universe, and I think the dataset is pretty static, so fetching everything and building up this structure is something we could do periodically. The benefit we'd get here is that we'd never encounter a search query where we don't have the results pre-fetched. We wouldn't need to do slow calls to SWAPI synchronously.

One unrelated performance change - right now the SWAPI calls are made in sequence. I think we could rework these to be made in parallel. This would give us performance improvement anywhere we do SWAPI search calls sequentially, regardless of the architecture. 