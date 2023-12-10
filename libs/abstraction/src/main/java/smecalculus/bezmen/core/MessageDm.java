package smecalculus.bezmen.core;

import lombok.Builder;
import lombok.NonNull;

public class MessageDm {
    @Builder
    public record RegistrationRequest(@NonNull String externalId) {}

    @Builder
    public record RegistrationResponse(@NonNull String externalId) {}

    @Builder
    public record PreviewRequest(@NonNull String externalId) {}

    @Builder
    public record PreviewResponse(@NonNull String externalId) {}
}
