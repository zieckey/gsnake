package freader

import (
	"flag"
)

var dir *string = flag.String("file_path", "e:/1/1", "The dir of the data file which we need to process")
var statusFile *string = flag.String("status", "e:/1/status.txt", "The status file which holds the processing status")
var priorityLevel *int = flag.Int("priority_level", 0, "The max priority level of the file handler. 0 means that it don't has any priorty")
var filePattern *string = flag.String("file_pattern", "inc_*.gz", "The pattern of the name which we need to process")
var reader_type *string = flag.String("reader_type", "PTailReader", "The type of the file reader, options supported now : GzipReader, PTailReader")
