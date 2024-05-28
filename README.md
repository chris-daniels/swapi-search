# Star Wars Character Search CLI

A quick integration with https://swapi.dev/api to search characters. 

### Requirements
Go version 1.22.0. 

### To run
```
$ go run .
```

### To test
```
$ go test .
```

### My approach
From what I could see, SWAPI only allows pagination as query parameters to the API. So in the worst case, we need to search through all pages of all SWAPI endpoints. 

I also assumed that the Star Wars canon is pretty static - a new Star Wars film likely won't come out during the lifetime of the CLI. 

The set of data working with here is pretty small, too. We can easily hold the names of characters/vehicles,films in-memory. 

So with all this in mind, I decided to simply crawl through all the pages of all of the SWAPI endpoints when the application starts and holding those values in-memory. The `searchTermCache` is a map of search terms to a set of character ids, and `peopleCache` maps those character ids to a `Person` pointer. There's some upfront cost to load all of this data, and it isn't resilient to SWAPI outages, but it seemed like a fair approach for this exercise. 

### Roadmap
My time was limited here so far, so I have a few items I'd like to try to get to. 
 - Improved testing
   - We should mock the SWAPI API when we test the `FetchSwapiData` function. 
   - We should test the API layer's `FetchEntities` function. 
 - I think it would be fun to try to do some fuzzy search here. 


### Bringing this to production
To bring this app to production, we would need to make a few changes. 

We could distribute this CLI tool, and people could use it that way, but eventually we might want something prettier. Also, having every individual user wait around for 6-7 seconds while they fetch the whole SWAPI API's data on startup doesn't make a ton of sense. If we centralize some stuff here, we can get a much faster experience for most users. 

So to move forward, I'd imagine we would have a client/server model in place here. The server would handle fetching data from the SWAPI and providing a search API, and clients would then handle taking user input and making search requests. This could still be a CLI tool, but we could also consider building a web front-end. 

I actually think my current approach, where Star Wars characters are held in memory, is pretty feasible in production. Whenever an instance of the service starts, we fetch the data, and every instance of the service can hold all this in memory. The scale of this data is small enough that the inpact on service startup time and memory consumed aren't too significant. 

If we wanted to pull from a broader dataset, though, this approach would fall apart. Suppose we wanted to pull data from all movies in existence. This set of data could grow large enough that it becomes expensive to hold everything in-memory. Even if we can fit everything on one machine, it's redundant to hold duplicated data across every machine. It would also start to take a load time to fetch all relevant data. Then we might want to reconsider our architecture and add a persistence layer. This could be a relational database where relationships between entities are managed by the RDMS.

Entries in our database should have a TTL, so we can fetch updates hourly/daily, whatever we decide is appropriate. We could set up an async task that writes a new version of everything, and then we update a "version" variable to point to the updated records. This way we can refresh our data without any impact on response time. 

With these changes, we'd have a super lightweight client that passes search terms along to the API. These requests are handled by machines that can be scaled horizontally, and these instances of our service will all share a persistence layer that has relatively up-to-date records. 

The only final thing I'd consider is standing up a task that polls the SWAPI endpoint and looks for updates. We could get `count` values off the response to detect when new data becomes available.  Then we can time our refreshes to fetch data as soon as possible. 