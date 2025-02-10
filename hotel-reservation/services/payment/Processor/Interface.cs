namespace Payment.Processor
{
    public interface IProcessor
    {
        Task<Models.Payment> Charge(Models.Payment payment);
    }
}
