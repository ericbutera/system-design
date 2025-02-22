namespace Payment.Processor
{
    public class FakeProcessor : IProcessor
    {
        public async Task<Models.Payment> Charge(Models.Payment payment)
        {
            await Task.Delay(300);
            payment.TransactionId = Guid.NewGuid().ToString();
            return payment;
        }
    }
}
