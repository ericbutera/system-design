"use client";

export default function HotelsList({
  hotels,
}: {
  hotels: { id: string; name: string; location: string }[];
}) {
  return (
    <div className="p-4">
      <h1 className="text-xl font-bold">Hotels</h1>
      <ul>
        {hotels.map((hotel) => (
          <li key={hotel.id} className="p-2 border-b">
            {hotel.name} - {hotel.location}
          </li>
        ))}
      </ul>
    </div>
  );
}
