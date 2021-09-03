.PHONY: protoc

protoc:
	 protoc -I=api/ --go_out=plugins=grpc:./ ./api/tf.proto