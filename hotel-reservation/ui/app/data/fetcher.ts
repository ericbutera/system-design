import createApolloClient from "@/app/lib/apolloClient";
import { gql } from "@apollo/client";

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

export async function fetchHotel(hotelId: string): Promise<HotelById> {
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
      id: Number(hotelId),
    },
  });
  return data.hotelById;
}

export interface Hotel {
  id: string;
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
