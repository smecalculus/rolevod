package smecalculus.bezmen.validation;

public interface EdgeValidator {
    <T> void validate(T object, Class<?>... groups);
}
