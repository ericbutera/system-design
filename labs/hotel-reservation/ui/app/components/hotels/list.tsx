import { fetchHotels } from "@/app/data/fetcher";
import Link from "next/link";

export default async function HotelsList() {
  const hotels = await fetchHotels();
  return (
    <div className="p-4">
      <ul>
        {hotels.map((hotel) => (
          <li key={hotel.id} className="p-2 border-b">
            <Link href={`/hotels/${hotel.id}`}>
              {hotel.name}: {hotel.location}
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
}
