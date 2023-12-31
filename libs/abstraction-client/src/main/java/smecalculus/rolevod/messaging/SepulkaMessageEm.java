package smecalculus.rolevod.messaging;

import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Data;

public abstract class SepulkaMessageEm {
    @Data
    public static class RegistrationRequest {
        @NotNull
        @Size(min = 1, max = 64)
        String externalId;
    }

    @Data
    public static class RegistrationResponse {
        @NotNull
        @Size(min = 1, max = 64)
        String externalId;
    }
}
