const { ApolloGateway, IntrospectAndCompose } = require("@apollo/gateway");
const { ApolloServer } = require("@apollo/server");
const express = require("express");
const { expressMiddleware } = require("@apollo/server/express4");

const hotelServiceUrl = process.env.HOTEL_SERVICE_URL || "http://hotel/graphql"; // Default if env var isn't set

const app = express();
app.use(express.json());

const gateway = new ApolloGateway({
  supergraphSdl: new IntrospectAndCompose({
    subgraphs: [
      { name: "hotel", url: hotelServiceUrl },
      // { name: 'reservation', url: 'http://reservation/graphql' },
      // { name: 'payment', url: 'http://payment/graphql' }
    ],
  }),
});

const server = new ApolloServer({
  gateway,
  introspection: true, // turn off in production
  plugins: [
    {
      requestDidStart() {
        return {
          didResolveOperation(context) {
            console.log("Resolved operation:", context.operationName);
          },
          willSendResponse(context) {
            console.log("Response:", context.response);
            if (context.errors) {
              console.error("Errors:", context.errors);
            }
          },
        };
      },
    },
  ],
});

// Apply Apollo middleware
server.start().then(() => {
  app.use("/graphql", expressMiddleware(server));

  app.listen(4000, () => {
    console.log("🚀 Apollo Gateway ready at http://localhost:4000/graphql");
  });
});

app.use((req, res, next) => {
  console.log("Incoming request:", req.body);
  if (!req.body.query) {
    console.error("Missing query in request body");
  }
  next();
});
