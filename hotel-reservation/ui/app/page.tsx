import Link from "next/link";
import { fetchHotels, Hotel } from "@/app/data/fetcher";

export default async function HomePage() {
  const hotels = await fetchHotels();
  return (
    <div className="p-4">
      <h1 className="text-xl font-bold">Hotels</h1>
      <ul>
        {hotels.map((hotel: Hotel) => (
          <li key={hotel.id} className="p-2 border-b">
            <Link
              href={`/hotels/${hotel.id}`}
              className="text-blue-500 hover:underline"
            >
              {hotel.name} - {hotel.location}
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
}
