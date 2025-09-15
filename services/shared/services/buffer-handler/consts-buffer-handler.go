package bufferhandler

const (
	LONG_TEXT = `
This time we get a RouteGuide_RouteChatServer stream that, as in our client-side streaming example, can be used to read and write messages. However, this time we return values via our method’s stream while the client is still writing messages to their message stream.
The syntax for reading and writing here is very similar to our client-streaming method, except the server uses the stream’s Send() method rather than SendAndClose() because it’s writing multiple responses. Although each side will always get the other’s messages in the order they were written, both the client and server can read and write in any order — the streams operate completely independently.
Starting the server Once we’ve implemented all our methods, we also need to start up a gRPC server so that clients can actually use our service. The following snippet shows how we do this for our RouteGuide service:
As in the simple RPC, we pass the method a context and a request. However, instead of getting a response object back, we get back an instance of RouteGuide_ListFeaturesClient. The client can use the RouteGuide_ListFeaturesClient stream to read the server’s responses.
We use the RouteGuide_ListFeaturesClient’s Recv() method to repeatedly read in the server’s responses to a response protocol buffer object (in this case a Feature) until there are no more messages: the client needs to check the error err returned from Recv() after each call. If nil, the stream is still good and it can continue reading; if it’s io.EOF then the message stream has ended; otherwise there must be an RPC error, which is passed over through err.
Client-side streaming RPC The client-side streaming method RecordRoute is similar to the server-side method, except that we only pass the method a context and get a RouteGuide_RecordRouteClient stream back, which we can use to both write and read messages.`
)
