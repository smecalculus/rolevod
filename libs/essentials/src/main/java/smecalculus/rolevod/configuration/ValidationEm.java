package smecalculus.rolevod.configuration;

import lombok.Data;

public abstract class ValidationEm {
    @Data
    public static class ValidationProps {
        private String mode;
    }
}
