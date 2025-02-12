import { fetchHotel, Room } from "@/app/data/fetcher";
import Link from "next/link";

interface Params {
  id: string;
}

export default async function HotelDetailsPage({
  params,
}: {
  params: Promise<Params>;
}) {
  const hotelId = (await params).id;
  const hotel = await fetchHotel(Number(hotelId));

  return (
    <div className="p-4">
      <h1 className="text-xl font-bold">{hotel.name}</h1>
      <p className="text-gray-600">{hotel.location}</p>

      <h2 className="text-lg font-semibold mt-4">Rooms</h2>
      <ul className="list-disc pl-5">
        {hotel.rooms.map((room: Room) => (
          <li key={room.id}>Room {room.number}</li>
        ))}
      </ul>

      <Link
        href={`/hotels/${hotel.id}/book`}
        className="text-blue-500 mt-4 block"
      >
        Book Now
      </Link>
    </div>
  );
}
