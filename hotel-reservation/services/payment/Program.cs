var builder = WebApplication.CreateBuilder(args);

// https://chillicream.com/docs/hotchocolate/v13/api-reference/apollo-federation/
builder.Services
    .AddGraphQLServer()
    .AddQueryType<Payment.GraphQL.Query>()
    .AddApolloFederation()
    .AddMutationType<Payment.GraphQL.Mutation>();

var app = builder.Build();

// app.UseHttpsRedirection();

app.MapGraphQL();

app.Run();
