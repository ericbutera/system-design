import { ApolloClient, InMemoryCache } from "@apollo/client";

const createApolloClient = (url?: string) => {
  return new ApolloClient({
    uri: url ?? "http://gateway/graphql",
    cache: new InMemoryCache(),
  });
};

export default createApolloClient;
