namespace Payment.Service
{
    public interface IPaymentService
    {
        Task<Models.Payment> Charge(Models.Payment payment);
        Task<Models.Payment> Get(int id);
    }
}
