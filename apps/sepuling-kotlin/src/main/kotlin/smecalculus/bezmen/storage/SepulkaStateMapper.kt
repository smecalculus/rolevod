package smecalculus.bezmen.storage

import org.mapstruct.Mapper
import smecalculus.bezmen.core.SepulkaStateDm
import smecalculus.bezmen.mapping.EdgeMapper

@Mapper
interface SepulkaStateMapper : EdgeMapper {
    fun toEdge(state: SepulkaStateDm.AggregateRoot): SepulkaStateEm.AggregateRoot

    fun toDomain(state: SepulkaStateEm.AggregateRoot): SepulkaStateDm.AggregateRoot

    fun toEdge(state: SepulkaStateDm.Touch): SepulkaStateEm.Touch

    fun toDomain(state: SepulkaStateEm.Existence): SepulkaStateDm.Existence

    fun toEdge(state: SepulkaStateDm.Preview): SepulkaStateEm.Preview

    fun toDomain(state: SepulkaStateEm.Preview): SepulkaStateDm.Preview
}
