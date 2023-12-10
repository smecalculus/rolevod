package smecalculus.bezmen.core;

import static java.time.temporal.ChronoUnit.MICROS;

import java.time.LocalDateTime;
import java.util.UUID;
import smecalculus.bezmen.core.StateDm.AggregateState;
import smecalculus.bezmen.core.StateDm.ExistenceState;
import smecalculus.bezmen.core.StateDm.PreviewState;
import smecalculus.bezmen.core.StateDm.TouchState;

public class StateDmEg {
    public static AggregateState.Builder aggregateState() {
        return AggregateState.builder()
                .internalId(UUID.randomUUID())
                .externalId(UUID.randomUUID().toString())
                .revision(0)
                .createdAt(LocalDateTime.now().truncatedTo(MICROS))
                .updatedAt(LocalDateTime.now().truncatedTo(MICROS));
    }

    public static ExistenceState.Builder creationState() {
        return ExistenceState.builder().internalId(UUID.randomUUID());
    }

    public static PreviewState.Builder previewState(AggregateState state) {
        return PreviewState.builder().externalId(state.externalId()).createdAt(state.createdAt());
    }

    public static TouchState.Builder touchState(AggregateState state) {
        return TouchState.builder().revision(state.revision()).updatedAt(state.updatedAt());
    }
}
