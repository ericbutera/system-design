from typing import Dict, List, Optional

import strawberry
from models import HotelModel, RoomModel
from strawberry.federation import Schema


@strawberry.federation.type(keys=["id"])
class Room:
    id: int
    number: str
    hotel_id: int


@strawberry.federation.type(keys=["id"])
class Hotel:
    id: int
    name: str
    location: str

    @strawberry.field
    def rooms(self, info: strawberry.Info) -> List["Room"]:
        db = info.context["db"]
        room_models = db.query(RoomModel).filter(RoomModel.hotel_id == self.id).all()
        return [
            Room(id=r.id, number=r.number, hotel_id=r.hotel_id) for r in room_models
        ]


@strawberry.type
class Query:
    @strawberry.field
    def hotels(self, info: strawberry.Info) -> List[Hotel]:
        db = info.context["db"]
        hotels = db.query(HotelModel).all()
        return [Hotel(id=h.id, name=h.name, location=h.location) for h in hotels]

    @strawberry.field
    def hotel_by_id(self, info: strawberry.Info, id: int) -> Optional[Hotel]:
        db = info.context["db"]
        h = db.query(HotelModel).filter(HotelModel.id == id).first()
        if h:
            return Hotel(id=h.id, name=h.name, location=h.location)
        return None


schema = Schema(query=Query)
