package smecalculus.bezmen.core;

import static java.time.temporal.ChronoUnit.MICROS;

import java.time.LocalDateTime;
import java.util.UUID;
import smecalculus.bezmen.core.SepulkaStateDm.AggregateRoot;
import smecalculus.bezmen.core.SepulkaStateDm.Preview;
import smecalculus.bezmen.core.SepulkaStateDm.Touch;

public class SepulkaStateDmEg {
    public static SepulkaStateDm.AggregateRoot.Builder aggregateRoot() {
        return AggregateRoot.builder()
                .internalId(UUID.randomUUID())
                .externalId(UUID.randomUUID().toString())
                .revision(0)
                .createdAt(LocalDateTime.now().truncatedTo(MICROS))
                .updatedAt(LocalDateTime.now().truncatedTo(MICROS));
    }

    public static SepulkaStateDm.Existence.Builder existence() {
        return SepulkaStateDm.Existence.builder().internalId(UUID.randomUUID());
    }

    public static SepulkaStateDm.Preview.Builder preview(SepulkaStateDm.AggregateRoot state) {
        return Preview.builder().externalId(state.externalId()).createdAt(state.createdAt());
    }

    public static SepulkaStateDm.Touch.Builder touch(AggregateRoot state) {
        return Touch.builder().revision(state.revision()).updatedAt(state.updatedAt());
    }
}
