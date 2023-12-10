package smecalculus.bezmen.storage

import org.springframework.data.annotation.Id
import org.springframework.data.domain.Persistable
import org.springframework.data.relational.core.mapping.Column
import org.springframework.data.relational.core.mapping.Table
import java.time.LocalDateTime
import java.util.UUID

abstract class SepulkaStateEm {
    data class Existence(
        var internalId: UUID,
        var externalId: String,
    )

    data class Preview(
        var internalId: UUID,
        var externalId: String,
        var createdAt: LocalDateTime,
    )

    data class Touch(
        var revision: Int,
        var updatedAt: LocalDateTime,
    )

    @Table("sepulkas")
    data class AggregateRoot(
        @Id
        var internalId: UUID?,
        @Column
        var externalId: String?,
        @Column
        var revision: Int?,
        @Column
        var createdAt: LocalDateTime?,
        @Column
        var updatedAt: LocalDateTime?,
    ) : Persistable<UUID> {
        override fun getId(): UUID? {
            return internalId
        }

        override fun isNew(): Boolean {
            return true
        }
    }
}
