package smecalculus.bezmen.storage.springdata;

import java.util.Optional;
import java.util.UUID;
import lombok.NonNull;
import org.springframework.data.jdbc.repository.query.Modifying;
import org.springframework.data.jdbc.repository.query.Query;
import org.springframework.data.repository.CrudRepository;
import org.springframework.data.repository.query.Param;
import smecalculus.bezmen.storage.SepulkaStateEm;

public interface SepulkaRepository extends CrudRepository<SepulkaStateEm.AggregateRoot, UUID> {

    Optional<SepulkaStateEm.Existence> findByExternalId(@NonNull String externalId);

    Optional<SepulkaStateEm.Preview> findByInternalId(@NonNull UUID internalId);

    @Modifying
    @Query(
            """
            UPDATE sepulkas
            SET revision = revision + 1,
                updated_at = :#{#state.updatedAt}
            WHERE internal_id = :id
            AND revision = :#{#state.revision}
            """)
    int updateBy(@Param("id") UUID internalId, @Param("state") SepulkaStateEm.Touch state);
}
