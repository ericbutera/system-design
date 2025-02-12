"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import {
  createReservation,
  ReservationInput,
  HotelById,
} from "@/app/data/fetcher";

export function CreateReservation({ hotel }: { hotel: HotelById }) {
  const [checkInDate, setCheckInDate] = useState("");
  const [checkOutDate, setCheckOutDate] = useState("");
  const [roomTypeId, setRoomTypeId] = useState(hotel.rooms[0]?.id || 0);
  const [quantity, setQuantity] = useState(1);
  const [error, setError] = useState("");
  const router = useRouter();

  //const roomTypes = await fetchRoomTypes(hotel.id);
  const roomTypes = [
    { id: 1, name: "Single" },
    { id: 2, name: "Double" },
  ];

  const handleReservation =
    (hotelId: number) => async (event: React.MouseEvent<HTMLButtonElement>) => {
      event.preventDefault();

      if (
        !checkInDate ||
        !checkOutDate ||
        !roomTypeId ||
        !quantity ||
        !hotelId
      ) {
        setError("All fields are required.");
        return;
      }
      setError("");

      const reservation = new ReservationInput(
        checkInDate,
        checkOutDate,
        roomTypeId,
        quantity,
        hotelId
      );

      try {
        const res = await createReservation(reservation);
        if (res.reservation) router.push(`/reservations/${res.reservation.id}`);
        return;
      } catch (e) {
        setError("An error occurred while booking." + e);
        return;
      }
    };

  return (
    <form>
      {error && <div className="error">{error}</div>}
      <div>
        <label>
          Check In Date
          <input
            type="date"
            name="checkInDate"
            value={checkInDate}
            required
            onChange={(e) => setCheckInDate(e.target.value)}
          />
        </label>
      </div>
      <div>
        <label>
          Check Out Date
          <input
            type="date"
            name="checkOutDate"
            value={checkOutDate}
            required
            onChange={(e) => setCheckOutDate(e.target.value)}
          />
        </label>
      </div>
      <div>
        <label>
          Room Type
          <select
            name="roomTypeId"
            value={roomTypeId}
            required
            onChange={(e) => setRoomTypeId(Number(e.target.value))}
          >
            {roomTypes.map((roomType) => (
              <option key={roomType.id} value={roomType.name}>
                {roomType.name}
              </option>
            ))}
          </select>
        </label>
      </div>
      <div>
        <label>
          Quantity
          <input
            type="number"
            name="quantity"
            value={quantity}
            min={1}
            max={5} // TODO: max available rooms
            required
            onChange={(e) => setQuantity(Number(e.target.value))}
          />
        </label>
      </div>
      <div>
        <button
          className="bg-blue-500 text-white px-4 py-2 rounded"
          onClick={handleReservation(Number(hotel.id))}
        >
          Confirm Booking
        </button>
      </div>
    </form>
  );
}
