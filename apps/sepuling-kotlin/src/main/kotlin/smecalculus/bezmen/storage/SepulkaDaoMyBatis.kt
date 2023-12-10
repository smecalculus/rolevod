package smecalculus.bezmen.storage

import smecalculus.bezmen.core.SepulkaStateDm
import smecalculus.bezmen.storage.mybatis.SepulkaSqlMapper
import java.util.UUID

class SepulkaDaoMyBatis(
    private val stateMapper: SepulkaStateMapper,
    private val sqlMapper: SepulkaSqlMapper,
) : SepulkaDao {
    override fun add(state: SepulkaStateDm.AggregateRoot): SepulkaStateDm.AggregateRoot {
        val stateEdge = stateMapper.toEdge(state)
        sqlMapper.insert(stateEdge)
        return state
    }

    override fun getBy(externalId: String): SepulkaStateDm.Existence? {
        return sqlMapper.findByExternalId(externalId)?.let { stateMapper.toDomain(it) }
    }

    override fun getBy(internalId: UUID): SepulkaStateDm.Preview? {
        return sqlMapper.findByInternalId(internalId.toString())?.let { stateMapper.toDomain(it) }
    }

    override fun updateBy(
        internalId: UUID,
        state: SepulkaStateDm.Touch,
    ) {
        val stateEdge = stateMapper.toEdge(state)
        val matchedCount = sqlMapper.updateBy(stateEdge, internalId.toString())
        if (matchedCount == 0) {
            throw ContentionException()
        }
    }
}
