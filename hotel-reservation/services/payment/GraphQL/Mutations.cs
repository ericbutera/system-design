namespace Payment.GraphQL
{
    public class Mutation
    {
        private readonly List<Transaction> _transactions = new List<Transaction>();

        public Transaction Charge(Transaction transaction)
        {
            // validation
            if (transaction.Amount <= 0)
            {
                throw new ArgumentException("Amount must be greater than zero.");
            }

            // var paymentResult = _paymentGateway.Charge(amount, ...);

            transaction.Id = Guid.NewGuid().ToString(); // TODO: paymentresult.transactionId

            _transactions.Add(transaction);

            return transaction;
        }
    }
}