package smecalculus.rolevod.validation;

public interface EdgeValidator {
    <T> void validate(T object, Class<?>... groups);
}
