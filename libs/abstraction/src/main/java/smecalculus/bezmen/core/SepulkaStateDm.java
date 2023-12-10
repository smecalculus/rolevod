package smecalculus.bezmen.core;

import java.time.LocalDateTime;
import java.util.UUID;
import lombok.Builder;
import lombok.NonNull;

public class SepulkaStateDm {
    @Builder
    public record Existence(@NonNull UUID internalId) {}

    @Builder
    public record Preview(@NonNull String externalId, @NonNull LocalDateTime createdAt) {}

    @Builder
    public record Touch(@NonNull Integer revision, @NonNull LocalDateTime updatedAt) {}

    @Builder
    public record AggregateRoot(
            @NonNull UUID internalId,
            @NonNull String externalId,
            @NonNull Integer revision,
            @NonNull LocalDateTime createdAt,
            @NonNull LocalDateTime updatedAt) {}
}
