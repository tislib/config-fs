package main

import (
	"github.com/thatisuday/commando"
)

func main() {
	backend := new(Backend)
	backend.init()

	// configure commando
	commando.
		SetExecutableName("config-fs").
		SetVersion("1.0.0").
		SetDescription("This tool lists the contents of a directory in tree-like format.\nIt can also display information about files and folders like size, permission and ownership.")

	commando.
		Register("read").
		SetShortDescription("Reads mongodb and stores data in fs").
		AddArgument("location", "location", "").                                                                     // required
		AddFlag("database,d", "which database to connect", commando.String, nil).                                    // required
		AddFlag("connection", "connection string to connect mongodb", commando.String, "mongodb://127.0.0.1:27017"). // required
		AddFlag("collections,c", "collections, collections filter by regex", commando.String, ".*").                 // required
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			location := args["location"].Value
			database := flags["database"].Value.(string)
			collection := flags["collections"].Value.(string)
			connection := flags["connection"].Value.(string)

			backend.runReadOperation(database, collection, connection, location)
		})

	commando.
		Register("write").
		SetShortDescription("Reads mongodb and stores data in fs").
		AddArgument("location", "location", "").                                                                     // required
		AddFlag("database,d", "which database to connect", commando.String, nil).                                    // required
		AddFlag("connection", "connection string to connect mongodb", commando.String, "mongodb://127.0.0.1:27017"). // required
		AddFlag("collections,c", "collections, collections filter by regex", commando.String, ".*").                 // required
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			location := args["location"].Value
			database := flags["database"].Value.(string)
			collection := flags["collections"].Value.(string)
			connection := flags["connection"].Value.(string)

			backend.runWriteOperation(database, collection, connection, location)
		})

	// parse command-line arguments
	commando.Parse(nil)

}
