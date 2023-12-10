package smecalculus.bezmen.validation;

import static java.lang.String.format;
import static java.util.stream.Collectors.toSet;

import jakarta.validation.ConstraintValidator;
import jakarta.validation.ConstraintValidatorContext;
import java.util.Set;
import java.util.stream.Stream;

public class ValueOfEnumValidator implements ConstraintValidator<ValueOfEnum, String> {
    private Set<String> allowedValues;

    @Override
    public void initialize(ValueOfEnum constraintAnnotation) {
        allowedValues = Stream.of(constraintAnnotation.value().getEnumConstants())
                .map(Enum::name)
                .collect(toSet());
    }

    @Override
    public boolean isValid(String value, ConstraintValidatorContext context) {
        if (value == null) {
            return true;
        }

        if (!allowedValues.contains(value.toUpperCase())) {
            context.buildConstraintViolationWithTemplate(format("Allowed values: %s", allowedValues))
                    .addConstraintViolation();
            return false;
        }

        return true;
    }
}
