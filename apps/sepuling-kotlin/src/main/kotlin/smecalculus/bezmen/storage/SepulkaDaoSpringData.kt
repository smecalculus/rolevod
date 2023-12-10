package smecalculus.bezmen.storage

import smecalculus.bezmen.core.SepulkaStateDm
import smecalculus.bezmen.storage.springdata.SepulkaRepository
import java.util.UUID

class SepulkaDaoSpringData(
    private val mapper: SepulkaStateMapper,
    private val repository: SepulkaRepository,
) : SepulkaDao {
    override fun add(state: SepulkaStateDm.AggregateRoot): SepulkaStateDm.AggregateRoot {
        val stateEdge = repository.save(mapper.toEdge(state))
        return mapper.toDomain(stateEdge)
    }

    override fun getBy(externalId: String): SepulkaStateDm.Existence? {
        return repository.findByExternalId(externalId)?.let { mapper.toDomain(it) }
    }

    override fun getBy(internalId: UUID): SepulkaStateDm.Preview? {
        return repository.findByInternalId(internalId)?.let { mapper.toDomain(it) }
    }

    override fun updateBy(
        internalId: UUID,
        state: SepulkaStateDm.Touch,
    ) {
        val stateEdge = mapper.toEdge(state)
        val matchedCount = repository.updateBy(internalId, stateEdge)
        if (matchedCount == 0) {
            throw ContentionException()
        }
    }
}
