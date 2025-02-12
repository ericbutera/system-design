import { fetchReservations } from "@/app/data/fetcher";
import Link from "next/link";

export default async function ReservationPage() {
  const reservations = await fetchReservations();

  return (
    <div className="p-4">
      <table>
        <thead>
          <tr>
            <th>Created</th>
            <th>Check In</th>
            <th>Check Out</th>
            <td>Status</td>
            <th>View</th>
          </tr>
        </thead>
        <tbody>
          {reservations.map((reservation) => (
            <tr key={reservation.id}>
              <td>{reservation.createdAt}</td>
              <td>{reservation.checkIn}</td>
              <td>{reservation.checkOut}</td>
              <td>{reservation.status}</td>
              <td>
                <Link href={`/reservations/${reservation.id}`}>Details</Link>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
