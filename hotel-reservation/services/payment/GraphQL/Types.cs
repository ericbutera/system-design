using System.ComponentModel.DataAnnotations;

namespace Payment.GraphQL
{
    public class Transaction
    {
        [ID]
        [Key]
        public required string Id { get; set; }
        public required string Status { get; set; }
        public decimal Amount { get; set; }
        public DateTime Timestamp { get; set; }
        public required string ReservationId { get; set; }
    }
}