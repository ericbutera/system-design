import createApolloClient from "@/app/lib/apolloClient";
import { gql } from "@apollo/client";

export interface Room {
  id: number;
  number: string;
  hotelId: number;
}

export interface HotelById {
  id: number;
  name: string;
  location: string;
  rooms: Room[];
}

export async function fetchHotel(id: number): Promise<HotelById> {
  const client = createApolloClient();
  const { data } = await client.query({
    query: gql`
      query HotelById($id: Int!) {
        hotelById(id: $id) {
          id
          name
          location
          rooms {
            id
            number
            hotelId
          }
        }
      }
    `,
    variables: {
      id: id,
    },
  });
  return data.hotelById;
}

export interface Hotel {
  id: number;
  name: string;
  location: string;
}

export async function fetchHotels(): Promise<Hotel[]> {
  const client = createApolloClient();
  const { data } = await client.query({
    query: gql`
      query Hotels {
        hotels {
          id
          name
          location
        }
      }
    `,
  });
  return data.hotels;
}

export class ReservationInput {
  constructor(
    public checkInDate: string,
    public checkOutDate: string,
    public roomTypeId: number,
    public quantity: number,
    public hotelId: number
  ) {}
}

const CREATE_MUTATION = gql`
  mutation CreateReservation(
    $hotelId: Int!
    $checkInDate: String!
    $checkOutDate: String!
    $roomTypeId: Int!
    $quantity: Int!
  ) {
    createReservation(
      input: {
        checkInDate: $checkInDate
        checkOutDate: $checkOutDate
        roomTypeId: $roomTypeId
        quantity: $quantity
        hotelId: $hotelId
      }
    ) {
      id
      quantity
      checkIn
      checkOut
      status
      roomTypeId
      hotelId
      paymentId
      guestId
      roomType
      createdAt
    }
  }
`;

export interface Reservation {
  id: number;
  quantity: number;
  checkIn: string;
  checkOut: string;
  status: string;
  roomTypeId: number;
  hotelId: number;
  paymentId: number;
  guestId: number;
  roomType: string;
  createdAt: string;
}

export interface CreateResult {
  reservation?: Reservation;
  errors?: string;
}
export async function createReservation(
  input: ReservationInput
): Promise<CreateResult> {
  const client = createApolloClient("http://localhost:4000/graphql");
  const res = await client.mutate({
    mutation: CREATE_MUTATION,
    variables: { ...input },
  });
  return { reservation: res.data.createReservation };
}

export interface Reservation {
  id: number;
  quantity: number;
  checkIn: string;
  checkOut: string;
  status: string;
  roomTypeId: number;
  roomType: string;
  createdAt: string;
}

export async function fetchReservation(id: number): Promise<Reservation> {
  const client = createApolloClient();
  const { data } = await client.query({
    query: gql`
      query ViewReservation {
        viewReservation(id: "2") {
          id
          quantity
          checkIn
          checkOut
          status
          roomTypeId
          hotelId
          paymentId
          guestId
          roomType
        }
      }
    `,
    variables: {
      id: id.toString(),
    },
  });
  return data.viewReservation;
}
