import { CreateReservation } from "@/app/components/reservation/form";
import { fetchHotel } from "@/app/data/fetcher";

interface Params {
  id: string;
}

export default async function BookingPage({
  params,
}: {
  params: Promise<Params>;
}) {
  const hotelId = (await params).id;
  const hotel = await fetchHotel(Number(hotelId));

  return (
    <div className="p-4">
      <h1 className="text-xl font-bold">Booking for Hotel {hotel.name}</h1>

      <CreateReservation hotel={hotel}></CreateReservation>
    </div>
  );
}
