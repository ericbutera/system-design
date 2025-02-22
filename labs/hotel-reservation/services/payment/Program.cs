using Microsoft.EntityFrameworkCore;
using Payment.GraphQL;
using Payment.Processor;
using Payment.Service;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddDbContext<ApplicationDbContext>(options =>
    options.UseNpgsql(builder.Configuration.GetConnectionString("DefaultConnection")));

builder.Services
    .AddScoped<Mutation>()
    .AddScoped<IProcessor, FakeProcessor>()
    .AddScoped<IPaymentService, PaymentService>()
    .AddErrorFilter<GraphQLErrorFilter>()
    .AddLogging()
    .AddGraphQLServer()
    .AddQueryType<Query>()
    .AddApolloFederation() // https://chillicream.com/docs/hotchocolate/v13/api-reference/apollo-federation/
    .AddMutationType<Mutation>();

var app = builder.Build();
// app.UseHttpsRedirection();
app.MapGraphQL();
app.Run();
