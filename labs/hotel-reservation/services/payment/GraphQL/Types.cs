using System.ComponentModel.DataAnnotations;

namespace Payment.GraphQL
{
    public class Transaction
    {
        [ID]
        public required string Id { get; set; }
        public required string CorrelationId { get; set; } // recordID from our system
        public required string TransactionId { get; set; } // payment gateway ID
        public decimal Amount { get; set; }
        public DateTime CreatedAt { get; set; }
    }

    public class ChargeInput
    {
        public decimal Amount { get; set; }
        public required string CorrelationId { get; set; }
    }
}
