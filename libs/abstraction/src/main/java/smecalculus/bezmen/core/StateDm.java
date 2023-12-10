package smecalculus.bezmen.core;

import java.time.LocalDateTime;
import java.util.UUID;
import lombok.Builder;
import lombok.NonNull;

public class StateDm {
    @Builder
    public record ExistenceState(@NonNull UUID internalId) {}

    @Builder
    public record PreviewState(@NonNull String externalId, @NonNull LocalDateTime createdAt) {}

    @Builder
    public record TouchState(@NonNull Integer revision, @NonNull LocalDateTime updatedAt) {}

    @Builder
    public record AggregateState(
            @NonNull UUID internalId,
            @NonNull String externalId,
            @NonNull Integer revision,
            @NonNull LocalDateTime createdAt,
            @NonNull LocalDateTime updatedAt) {}
}
