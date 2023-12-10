package smecalculus.bezmen.storage.springdata

import org.springframework.data.jdbc.repository.query.Modifying
import org.springframework.data.jdbc.repository.query.Query
import org.springframework.data.repository.CrudRepository
import org.springframework.data.repository.query.Param
import smecalculus.bezmen.storage.SepulkaStateEm
import java.util.UUID

interface SepulkaRepository : CrudRepository<SepulkaStateEm.AggregateRoot, UUID> {
    fun findByExternalId(externalId: String): SepulkaStateEm.Existence?

    fun findByInternalId(internalId: UUID): SepulkaStateEm.Preview?

    @Modifying
    @Query(
        """
        UPDATE sepulkas
        SET revision = revision + 1,
            updated_at = :#{#state.updatedAt}
        WHERE internal_id = :id
        AND revision = :#{#state.revision}
        """,
    )
    fun updateBy(
        @Param("id") internalId: UUID,
        @Param("state") state: SepulkaStateEm.Touch,
    ): Int
}
