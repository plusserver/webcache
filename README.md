# webcache
Lightweight memory cache for go web servers that you can wrap around time consuming requests

# What does it do?
If you write a basic webserver with go and httpHandler() you might have some functions that take longer to return data. In a case where these functions are requested more often than your server can reply them, it will be like a DoS. This package helps you to simply wrap the content of those "expensive" functions in a memory cache.
During the update of the content the client that triggered the update is blocked and has to wait. All other clients will receive the old data without waiting. So only one client is waiting for the data, while all others are served from the cache.

# Usage
If your code is like this

	func main() {
		...
		http.HandleFunc("/myapi/", apiHandler)
		...
	}
	// request handler with an expensive backend function
	func apiHandler(w http.ResponseWriter, r *http.Request) {
		...
		result := complexBackendFunction() // takes long time to come back
		w.Write(result)
	}
	
you have to add or modify some lines, so it looks like this

	func main() {
		...
		apiCache = webcache.NewCachedPage(time.Second * 90) // cache for 90 seconds
		http.HandleFunc("/myapi/", apiHandler)
		...
	}
	// request handler with an expensive backend function
	func apiHandler(w http.ResponseWriter, r *http.Request) {
		...
		if !apiCache.Valid() {
			if apiCache.StartUpdate() == nil {
				apiCache.Write(complexBackendFunction())
				apiCache.EndUpdate()
			}
		}
		w.Write(apiCache.Get())
	}

If you need to invalidate the cache manually, call the `.Clear()` method. A full example server implementation can be found [here](https://github.com/Nexinto/webcache/tree/master/example%20server/main.go)
