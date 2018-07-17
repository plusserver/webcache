# webcache
Lightweight memory cache for go web servers that you can wrap around time consuming requests

# What does it do?
If you write a basic webserver with go a the httpHandler you might have some functions that take longer to return data. If these functions are requested quicker than your server can reply you end up in a DoS szenario. This package helps you to simply wrap thos functions in a memory cache. 


# Usage
If your code is like this

	func main() {
		...
		http.HandleFunc("/myapi/", apiHandler)
		...
	}
	// simple handler with a expensive backend function
	func apiHandler(w http.ResponseWriter, r *http.Request) {
		...
		result := complexBackendFunction() // takes long time to come back
		w.Write(result)
	}
	
you have to add or modify some lines, so it looks like this

	func main() {
		...
		apiCache = webcache.NewCachedPage(time.Second * 10)
		http.HandleFunc("/myapi/", apiHandler)
		...
	}
	// simple handler with a expensive backend function
	func apiHandler(w http.ResponseWriter, r *http.Request) {
		...
		if !apiCache.Valid() {
			if apiCache.StartUpdate() == nil {
				result := complexBackendFunction() // takes long time to come back
				apiCache.EndUpdate()
			}
		}
		w.Write(apiCache.Get())
	}

