package aws

// import (
// 	"fmt"

// 	"github.com/scottd018/policy-gen/internal/pkg/processor"
// )

// // NewProcessor instantiates a new instance of a markerProcessor
// // object.
// func NewProcessor(config *processor.Config) (*processor.Processor, error) {
// 	markerProcessor, err := processor.NewProcessor(config, MarkerDefinition(), Marker{})
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to create marker processor - %w", err)
// 	}

// 	// add the generator
// 	markerProcessor.Config.FileGenerator = &PolicyFileGenerator{
// 		Directory: config.OutputDirectory,
// 	}

// 	return markerProcessor, nil
// }
