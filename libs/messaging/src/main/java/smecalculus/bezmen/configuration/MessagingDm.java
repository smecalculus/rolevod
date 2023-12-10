package smecalculus.bezmen.configuration;

import java.util.Set;
import lombok.Builder;
import lombok.NonNull;

public abstract class MessagingDm {
    public enum MappingMode {
        SPRING_MVC,
        SPRING_MANAGEMENT
    }

    public enum ProtocolMode {
        HTTP,
        JMX
    }

    @Builder
    public record MessagingProps(@NonNull ProtocolProps protocolProps, @NonNull MappingProps mappingProps) {}

    @Builder
    public record ProtocolProps(@NonNull Set<ProtocolMode> protocolModes) {}

    @Builder
    public record MappingProps(@NonNull Set<MappingMode> mappingModes) {}
}
