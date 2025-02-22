# Design a Parking Garage

Low level design for a parking garage.

## Requirements

- Parking lot has multiple levels
- Each level has spots
- Parking lot can park motorcycles, cars, and buses
- A motorcycle can park in any spot
- A car can park in a single spot
- A bus can park in five contiguous spots

## Class Diagram

```mermaid
classDiagram
  class Garage {
    string +Name
    List~Level~ +Levels
    +Reserve(Vehicle) Tuple~Reservation, error~
    +Checkin(Vehicle) Tuple~Ticket, error~
    +Checkout(Ticket) Tuple~Receipt, error~
  }

  class Level {
    List~Spot~ +Spots
  }

  class Spot {
    Ticket +Ticket
  }

  class Ticket {
    time.Time +Checkin
    time.Time +Checkout
  }

  class Vehicle {
    <<abstract>>
    +Size() int
    +Ticket() Ticket
  }

  Vehicle <|-- Car
  Vehicle <|-- Motorcycle
  Vehicle <|-- Bus

  class Reservation {
    Level +Level
    List~Spot~ +Spots
  }
```

## Entity Relationships

```mermaid
erDiagram
  Garage
  Level
  Spot
  Car
  Ticket
  Receipt

  Garage ||--o| Level : has
  Level ||--o|Spot : has
  Spot ||--o|Ticket : has

  Car ||--o|Ticket : has
  Receipt ||--o|Garage : has

  Reservation ||--o|Level : has
  Reservation ||--o|Spot : has

  Bus ||--o|Vehicle : is
  Motorcycle ||--o|Vehicle : is
  Car ||--o|Vehicle : is
```

## Workflow

```mermaid
flowchart LR
  Vehicle --> Checkin
  Checkin --> Reservation{Availability?}
  Reservation --> Yes
  Reservation --> No
  No --> Checkin
  Yes --> |Ticket| Park
  Park --> Checkout
  Checkout --> |Receipt| Exit
```

## Out of Scope

- Gates
- Vehicle methods: Park, Drive, etc
- Custom Rates (Garage, vehicle type, etc.)
- Payment
  - PaymentGateway
  - Types: CC, Parking Pass
- Reporting
  - Currently available
  - Currently reserved
  - Breakdowns by vehicle type
  - Sales
  - Availability over time
