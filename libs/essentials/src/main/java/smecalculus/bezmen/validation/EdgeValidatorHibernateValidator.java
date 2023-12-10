package smecalculus.bezmen.validation;

import jakarta.validation.ConstraintViolationException;
import jakarta.validation.Validator;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;

@RequiredArgsConstructor
public class EdgeValidatorHibernateValidator implements EdgeValidator {

    @NonNull
    private Validator validator;

    @Override
    public <T> void validate(T object, Class<?>... groups) {
        var violations = validator.validate(object, groups);
        if (!violations.isEmpty()) {
            throw new ConstraintViolationException(violations);
        }
    }
}
