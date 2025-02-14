import {
  ApolloClient,
  InMemoryCache,
  HttpLink,
  ApolloLink,
} from "@apollo/client";
import { onError } from "@apollo/client/link/error";

const errorLink = onError(({ graphQLErrors, networkError }) => {
  if (graphQLErrors) {
    graphQLErrors.forEach(({ message, locations, path }) => {
      console.error(
        `GraphQL error ${message}, Location: ${locations}, Path: ${path}`
      );
    });
  }
  if (networkError) {
    console.error(`Network error: ${networkError}`);
  }
});

export default function createApolloClient(url?: string) {
  const uri = url ?? "http://gateway/graphql";
  const httpLink = new HttpLink({ uri: uri });
  const link = ApolloLink.from([errorLink, httpLink]);
  return new ApolloClient({
    link,
    cache: new InMemoryCache(),
  });
}
