package ddd

type (
        AggregateNamer interface {
                AggregateName() string                
        }

        Eventer interface {
                AddEvent(string, EventPayload, ...EventOption) 
                Events() []AggregateEvent
                ClearEvents()
        }


)
