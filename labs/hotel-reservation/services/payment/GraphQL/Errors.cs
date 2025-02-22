namespace Payment.GraphQL
{
    // https://josiahmortenson.dev/blog/2020-06-05-hotchocolate-graphql-errors
    public class GraphQLErrorFilter : IErrorFilter
    {
        public IError OnError(IError error)
        {
            return error.Exception is null ? error : error.WithMessage(error.Exception.Message);
        }
    }
}
