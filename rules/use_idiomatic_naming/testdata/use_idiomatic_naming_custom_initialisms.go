package fixtures

// With extraInitialisms: ["GRPC", "AMQP"]

// Invalid: Custom initialism not capitalized
type GrpcServer struct{} // MATCH /type GrpcServer should be GRPCServer/

type AmqpConnection struct{} // MATCH /type AmqpConnection should be AMQPConnection/

// Invalid: Default initialism still applies
type HttpClient struct{} // MATCH /type HttpClient should be HTTPClient/

// Valid: Correctly capitalized custom initialisms
type GRPCServer struct{}

type AMQPConnection struct{}

// Valid: Correctly capitalized default initialism
type HTTPClient struct{}
