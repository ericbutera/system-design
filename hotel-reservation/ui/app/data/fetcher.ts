import createApolloClient from "@/app/lib/apolloClient";
import { ApolloError, gql } from "@apollo/client";

export interface Room {
  id: string;
  number: string;
  hotelId: string;
}

export interface HotelById {
  id: string;
  name: string;
  location: string;
  rooms: Room[];
}

const FETCH_HOTEL = gql`
  query HotelById($id: ID!) {
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
`;
export async function fetchHotel(id: string): Promise<HotelById> {
  const client = createApolloClient();
  try {
    const { data } = await client.query({
      query: FETCH_HOTEL,
      variables: { id: id },
    });
    return data.hotelById;
  } catch (e) {
    console.error(e);
    if (e instanceof ApolloError) {
      console.error("gql errors %o", e.graphQLErrors);
      console.error("network errors %o", e.networkError);
    }
    throw e;
  }
}

export interface Hotel {
  id: string;
  name: string;
  location: string;
}

const FETCH_HOTELS = gql`
  query Hotels {
    hotels {
      id
      name
      location
    }
  }
`;
export async function fetchHotels(): Promise<Hotel[]> {
  const client = createApolloClient();
  const { data } = await client.query({
    query: FETCH_HOTELS,
  });
  return data.hotels;
}

export class ReservationInput {
  constructor(
    public checkInDate: string,
    public checkOutDate: string,
    public roomTypeId: string,
    public quantity: number,
    public hotelId: string
  ) {}
}

const CREATE_MUTATION = gql`
  mutation CreateReservation(
    $hotelId: ID!
    $checkInDate: String!
    $checkOutDate: String!
    $roomTypeId: ID!
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
  id: string;
  quantity: number;
  checkIn: string;
  checkOut: string;
  status: string;
  roomTypeId: string;
  hotelId: string;
  paymentId: string;
  guestId: string;
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
  id: string;
  quantity: number;
  checkIn: string;
  checkOut: string;
  status: string;
  roomTypeId: string;
  roomType: string;
  createdAt: string;
}

const VIEW_RESERVATION = gql`
  query ViewReservation($id: ID!) {
    viewReservation(id: $id) {
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
`;
export async function fetchReservation(id: string): Promise<Reservation> {
  const client = createApolloClient();
  const { data } = await client.query({
    query: VIEW_RESERVATION,
    variables: { id: id },
  });
  return data.viewReservation;
}

const VIEW_RESERVATIONS = gql`
  query ViewReservations {
    viewReservations {
      id
      quantity
      checkIn
      checkOut
      status
      roomTypeId
      paymentId
      guestId
      roomType
      createdAt
      hotelId
    }
  }
`;
export async function fetchReservations(): Promise<Reservation[]> {
  const client = createApolloClient();
  const { data } = await client.query({
    query: VIEW_RESERVATIONS,
  });
  return data.viewReservations;
}
