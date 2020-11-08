package main

type archiveProvider interface {
	// submit is responsible for send request to provider(archive.org).
	// param:
	//		- url: user input url
	// return:
	//		- html: text that represent archive page
	//		- err: error
	submit(url string) (html string, err error)

	// analysis is response for get the unique archive result url. If your provider doesn't require this, just write
	// an empty responsible that will return html.
	// param:
	//			- html: above result
	// return:
	//			- unique: identifier
	//			- err: error
	analysis(html string) (unique string, err error)

	// status is responsible for refresh the newest archive status because archive will take some time.
	// param:
	//			- unique: the unique stuff to represent the archive page
	// return:
	//			- result: archive progress
	//			- err: error
	status(unique string) (result string, err error)
}
