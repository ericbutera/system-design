const {
  ApolloGateway,
  IntrospectAndCompose,
  RemoteGraphQLDataSource,
} = require("@apollo/gateway");
const { startStandaloneServer } = require("@apollo/server/standalone");
const { ApolloServer } = require("@apollo/server");
const express = require("express");

const hotel = process.env.HOTEL_URL || "http://hotel/graphql";
const reservation = process.env.RESERVATION_URL || "http://reservation/query";
const payment = process.env.PAYMENT_URL || "http://payment/graphql/";
// const search = process.env.SEARCH_URL || "http://search/graphql";

const app = express();
app.use(express.json());

class AuthenticatedDataSource extends RemoteGraphQLDataSource {
  willSendRequest({ request, context }) {
    // in a real system pass the JWT token to the subgraph
    request.http.headers.set("user-id", context.userId);
  }
}

const supergraphSdl = new IntrospectAndCompose({
  subgraphs: [
    { name: "hotel", url: hotel },
    { name: "reservation", url: reservation },
    { name: "payment", url: payment },
    // { name: "search", url: search },
  ],
});
const gateway = new ApolloGateway({
  supergraphSdl: supergraphSdl,
  buildService({ name, url }) {
    return new AuthenticatedDataSource({ url });
  },
  serviceHealthCheck: true,
  pollIntervalInMs: 30_000, // 30 seconds
});

const reqPlugin = {
  requestDidStart() {
    return {
      didResolveOperation(context) {
        console.log("Resolved operation:", context.operationName);
      },
      willSendResponse(context) {
        if (context.errors) console.error("Errors:", context.errors);
      },
    };
  },
};
const server = new ApolloServer({
  gateway,
  introspection: true, // turn off in production
  plugins: [reqPlugin],
});

// https://www.apollographql.com/docs/apollo-server/using-federation/apollo-gateway-setup#advanced-usage
const getUserId = (token) => {
  return 1;
};

startStandaloneServer(server, {
  context: ({ req }) => {
    const token = req.headers.authorization || ""; // get JWT token from header
    const userId = getUserId(token);
    return { userId };
  },
});
console.log(`Server ready`);
