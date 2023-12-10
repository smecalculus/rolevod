package smecalculus.bezmen.storage;

import org.mapstruct.Mapper;
import smecalculus.bezmen.core.StateDm;
import smecalculus.bezmen.mapping.EdgeMapper;

@Mapper
public interface SepulkaStateMapper extends EdgeMapper {
    StateEm.AggregateState toEdge(StateDm.AggregateState state);

    StateDm.AggregateState toDomain(StateEm.AggregateState state);

    StateEm.TouchState toEdge(StateDm.TouchState state);

    StateDm.ExistenceState toDomain(StateEm.ExistenceState state);

    StateEm.PreviewState toEdge(StateDm.PreviewState state);

    StateDm.PreviewState toDomain(StateEm.PreviewState state);
}
