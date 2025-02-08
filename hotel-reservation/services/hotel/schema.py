from typing import List

import strawberry
from strawberry.federation import Schema


@strawberry.federation.type(keys=["id"])
class Hotel:
    id: int
    name: str
    location: str
    rating: float


hotels = [
    Hotel(id=1, name="Grand Plaza", location="New York", rating=4.5),
    Hotel(id=2, name="Ocean View", location="Miami", rating=4.7),
]


@strawberry.type
class Query:
    hotels: List[Hotel] = strawberry.field(default_factory=lambda: hotels)


schema = Schema(query=Query)
