package smecalculus.bezmen.storage;

import java.util.Optional;
import java.util.UUID;
import smecalculus.bezmen.core.StateDm.AggregateState;
import smecalculus.bezmen.core.StateDm.ExistenceState;
import smecalculus.bezmen.core.StateDm.PreviewState;
import smecalculus.bezmen.core.StateDm.TouchState;

/**
 * Port: server side
 */
public interface SepulkaDao {
    AggregateState add(AggregateState state);

    Optional<ExistenceState> getBy(String externalId);

    Optional<PreviewState> getBy(UUID internalId);

    void updateBy(TouchState state, UUID internalId);
}
