from typing import List

from sqlalchemy import Column, Float, ForeignKey, Integer, String, create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import relationship, sessionmaker
from strawberry.federation import Schema

DATABASE_URL = "postgresql://postgres:password@pg-postgresql:5432/postgres"  # TODO env

Base = declarative_base()
engine = create_engine(DATABASE_URL, echo=True)
SessionLocal = sessionmaker(bind=engine, autoflush=False, expire_on_commit=False)


# Define SQLAlchemy models
class HotelModel(Base):
    __tablename__ = "hotels"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String, nullable=False)
    location = Column(String, nullable=False)
    rooms = relationship("RoomModel", back_populates="hotel")


class RoomModel(Base):
    __tablename__ = "rooms"

    id = Column(Integer, primary_key=True, index=True)
    hotel_id = Column(Integer, ForeignKey("hotels.id"))
    number = Column(String, nullable=False)
    hotel = relationship("HotelModel", back_populates="rooms")
