package smecalculus.bezmen.storage;

import java.util.Optional;
import java.util.UUID;
import smecalculus.bezmen.core.SepulkaStateDm;

/**
 * Port: server side
 */
public interface SepulkaDao {
    SepulkaStateDm.AggregateRoot add(SepulkaStateDm.AggregateRoot state);

    Optional<SepulkaStateDm.Existence> getBy(String externalId);

    Optional<SepulkaStateDm.Preview> getBy(UUID internalId);

    void updateBy(UUID internalId, SepulkaStateDm.Touch state);
}
