namespace Payment.GraphQL
{
    public class Query
    {
        private readonly ApplicationDbContext _context;

        public Query(ApplicationDbContext context)
        {
            _context = context;
        }

        public Transaction GetTransaction(Guid id)
        {
            var payment = _context.Payments.Find(id);
            if (payment == null)
            {
                throw new InvalidOperationException("Payment not found.");
            }

            return new Transaction
            {
                Id = payment.Id.ToString(),
                Amount = payment.Amount,
                CorrelationId = payment.CorrelationId,
                TransactionId = payment.TransactionId ?? string.Empty,
                CreatedAt = payment.CreatedAt,
            };
        }

    }
}
