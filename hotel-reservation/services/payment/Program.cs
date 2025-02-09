var builder = WebApplication.CreateBuilder(args);

builder.Services
    .AddGraphQLServer()
    .AddQueryType<Payment.GraphQL.Query>()
    .AddApolloFederation()
    .AddMutationType<Payment.GraphQL.Mutation>();

var app = builder.Build();

// app.UseHttpsRedirection();

app.MapGraphQL();

app.Run();
