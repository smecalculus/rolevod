package smecalculus.rolevod.storage;

import org.mapstruct.Mapper;
import smecalculus.rolevod.core.SepulkaStateDm;
import smecalculus.rolevod.mapping.EdgeMapper;

@Mapper
public interface SepulkaStateMapper extends EdgeMapper {
    SepulkaStateEm.AggregateRoot toEdge(SepulkaStateDm.AggregateRoot state);

    SepulkaStateDm.AggregateRoot toDomain(SepulkaStateEm.AggregateRoot state);

    SepulkaStateEm.Touch toEdge(SepulkaStateDm.Touch state);

    SepulkaStateDm.Existence toDomain(SepulkaStateEm.Existence state);

    SepulkaStateEm.Preview toEdge(SepulkaStateDm.Preview state);

    SepulkaStateDm.Preview toDomain(SepulkaStateEm.Preview state);
}
