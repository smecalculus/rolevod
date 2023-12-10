package smecalculus.bezmen.configuration;

import lombok.Builder;
import lombok.NonNull;

public abstract class ValidationDm {
    public enum ValidationMode {
        HIBERNATE_VALIDATOR
    }

    @Builder
    public record ValidationProps(@NonNull ValidationMode validationMode) {}
}
