package smecalculus.bezmen.storage

import smecalculus.bezmen.core.SepulkaStateDm
import java.util.UUID

/**
 * Port: server side
 */
interface SepulkaDao {
    fun add(state: SepulkaStateDm.AggregateRoot): SepulkaStateDm.AggregateRoot

    fun getBy(externalId: String): SepulkaStateDm.Existence?

    fun getBy(internalId: UUID): SepulkaStateDm.Preview?

    fun updateBy(
        internalId: UUID,
        state: SepulkaStateDm.Touch,
    )
}
