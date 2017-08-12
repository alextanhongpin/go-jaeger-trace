If you hit the following errors:

```bash
# github.com/uber/jaeger-client-go/thrift-gen/baggage
../../uber/jaeger-client-go/thrift-gen/baggage/baggagerestrictionmanager.go:155: cannot use baggageRestrictionManagerProcessorGetBaggageRestrictions literal (type *baggageRestrictionManagerProcessorGetBaggageRestrictions) as type thrift.TProcessorFunction in assignment:
	*baggageRestrictionManagerProcessorGetBaggageRestrictions does not implement thrift.TProcessorFunction (wrong type for Process method)
		have Process(int32, thrift.TProtocol, thrift.TProtocol) (bool, thrift.TException)
		want Process(context.Context, int32, thrift.TProtocol, thrift.TProtocol) (bool, thrift.TException)
```

Install glide:

$ brew install glide

Rebuild
```bash
go get -u github.com/uber/jaeger-client-go/
cd $GOPATH/src/github.com/uber/jaeger-client-go/
git submodule update --init --recursive
make install
```
http://localhost:16686/search