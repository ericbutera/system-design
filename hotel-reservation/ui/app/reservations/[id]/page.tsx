import { fetchReservation } from "@/app/data/fetcher";
import Link from "next/link";

interface Params {
  id: string;
}

export default async function ReservationPage({
  params,
}: {
  params: Promise<Params>;
}) {
  const id = (await params).id;
  const reservation = await fetchReservation(Number(id));

  return (
    <div className="p-4">
      <h2>Reservation {reservation.id}</h2>
      <p>Hotel: {reservation.hotelId}</p>
      <p>Room: {reservation.roomType}</p>
      <p>Check-in: {reservation.checkIn}</p>
      <p>Check-out: {reservation.checkOut}</p>
      <p>Quantity: {reservation.quantity}</p>
      <p>Status: {reservation.status}</p>
      <p>Created: {reservation.createdAt}</p>
      <Link href={`/hotels/${reservation.hotelId}`}>Hotel Details</Link>
    </div>
  );
}
