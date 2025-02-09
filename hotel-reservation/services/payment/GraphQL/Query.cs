namespace Payment.GraphQL
{
    public class Query
    {
        public Transaction GetTransaction(Guid id) =>
            new Transaction
            {
                Id = id.ToString(),
                Status = "Success",
                Amount = 100,
                Timestamp = DateTime.UtcNow,
                ReservationId = "123"
            };
    }
}