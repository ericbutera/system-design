import HotelsList from "./components/hotels/list";

export default async function HomePage() {
  return (
    <div className="p-4">
      <h1 className="text-xl font-bold">Hotels</h1>
      <HotelsList></HotelsList>
    </div>
  );
}
