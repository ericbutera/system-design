using Microsoft.EntityFrameworkCore;
using Payment.Service;

namespace Payment.GraphQL
{
    public class Mutation
    {
        private readonly IPaymentService _paymentService;

        public Mutation(IPaymentService paymentService)
        {
            _paymentService = paymentService;
        }

        public async Task<Transaction> Charge(ChargeInput charge)
        {
            var payment = new Models.Payment
            {
                Amount = charge.Amount,
                CorrelationId = charge.CorrelationId,
            };

            payment = await _paymentService.Charge(payment);

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
