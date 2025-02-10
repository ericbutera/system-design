
using Microsoft.EntityFrameworkCore;
using Npgsql;
using KeyNotFoundException = System.Collections.Generic.KeyNotFoundException;
using Payment.Processor;

namespace Payment.Service
{
    public class PaymentService : IPaymentService
    {

        private readonly ApplicationDbContext _context;
        private readonly ILogger<PaymentService> _logger;
        private readonly IProcessor _processor;

        public PaymentService(ApplicationDbContext context, ILogger<PaymentService> logger, IProcessor processor)
        {
            _context = context;
            _logger = logger;
            _processor = processor;
        }

        public async Task<Models.Payment> Charge(Models.Payment payment)
        {
            _logger.LogInformation("Charge correlationId: {CorrelationId} amount: {Amount}", payment.CorrelationId, payment.Amount);
            if (payment.Amount <= 0)
            {
                throw new ArgumentException("Amount must be greater than zero.");
            }
            else if (string.IsNullOrWhiteSpace(payment.CorrelationId))
            {
                throw new ArgumentException("CorrelationId must not be empty.");
            }

            payment = await _processor.Charge(payment);

            try
            {
                await _context.Payments.AddAsync(payment);
                payment.Id = await _context.SaveChangesAsync();
            }
            catch (DbUpdateException ex) when (ex.InnerException is PostgresException pgEx && pgEx.SqlState == "23505")
            {
                throw new DuplicateChargeException("Duplicate charge detected.");
            }
            return payment;
        }

        public async Task<Models.Payment> Get(int id)
        {
            var payment = await _context.Payments.FindAsync(id);
            if (payment == null)
            {
                throw new KeyNotFoundException($"Payment with id {id} not found.");
            }
            return payment;
        }
    }

    [Serializable]
    internal class DuplicateChargeException : Exception
    {
        public DuplicateChargeException()
        {
        }

        public DuplicateChargeException(string? message) : base(message)
        {
        }

        public DuplicateChargeException(string? message, Exception? innerException) : base(message, innerException)
        {
        }
    }
}
